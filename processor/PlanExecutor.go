package processor

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/request"

	"github.com/REAANDREW/telegraph"
)

//PlanExecutor ...
type PlanExecutor struct {
	Config    *config.Configuration
	Stats     *Statistics
	Bar       ProgressBar
	start     time.Time
	publisher telegraph.LinkedPublisher
}

func (instance *PlanExecutor) createPlan() Plan {
	plan := Plan{
		Name:     "Plan from urls in file",
		Workers:  instance.Config.Workers,
		WaitTime: instance.Config.WaitTime,
		Jobs:     []Job{},
	}

	job := Job{
		Name:  "Job for the urls in file",
		Steps: []Step{},
	}

	reader := request.NewRequestReader(instance.Config.FilePath)

	stream := request.NewSequentialRequestStream(reader)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}

	for stream.HasNext() {
		request, _ := stream.Next()
		step := Step{}

		action := &HTTPRequestExecutionAction{
			Client:  client,
			URL:     request.URL.String(),
			Method:  request.Method,
			Headers: request.Header,
		}

		step.Action = action
		job.Steps = append(job.Steps, step)
	}

	plan.Jobs = append(plan.Jobs, job)

	return plan
}

//ExecutionStarted ...
type ExecutionStarted struct{}

//ExecutionStopped ...
type ExecutionStopped struct{}

var resultHandlers = map[string]func(obj interface{}, statistics *Statistics){
	"http:request:error": func(obj interface{}, statistics *Statistics) {
		statistics.Request(obj.(error))
	},
	"http:response:error": func(obj interface{}, statistics *Statistics) {
		statistics.Request(obj.(error))
	},
	"http:request:bytes": func(obj interface{}, statistics *Statistics) {
		statistics.BytesSent(int64(obj.(int)))
	},
	"http:response:bytes": func(obj interface{}, statistics *Statistics) {
		statistics.BytesReceived(int64(obj.(int)))
	},
	"http:response:status": func(obj interface{}, statistics *Statistics) {
		statistics.Request(nil)
	},
	"action:duration": func(obj interface{}, statistics *Statistics) {
		statistics.ResponseTime(int64(obj.(time.Duration)))
	},
}

func (instance *PlanExecutor) executeStep(step Step) ExecutionResult {
	start := time.Now()
	executionResult := step.Action.Execute()
	duration := time.Since(start) / time.Millisecond
	executionResult["action:duration"] = duration
	for _, assertion := range step.Assertions {
		assertionResult := assertion.Assert(executionResult)
		executionResult[assertion.ResultKey()] = assertionResult
	}
	return executionResult
}

func (instance *PlanExecutor) workerExecuteJobs(jobs []Job) {
	for _, job := range jobs {
		func(talula Job) {
			var stepStream StepStream
			stepStream = CreateStepSequentialStream(talula.Steps)
			if instance.Config.Random {
				stepStream = CreateStepRandomStream(talula.Steps)
			}
			if instance.Config.WaitTime > time.Duration(0) {
				stepStream = CreateStepDelayStream(stepStream, instance.Config.WaitTime)
			}
			for stepStream.HasNext() {
				step := stepStream.Next()
				executionResult := instance.executeStep(step)
				instance.publisher.Publish(executionResult)
				if instance.Config.Duration > 0 && time.Since(instance.start) > instance.Config.Duration {
					break
				}
			}
		}(job)
		if instance.Config.Duration > 0 && time.Since(instance.start) < instance.Config.Duration {
			instance.executeJobs(jobs)
		} else {
			break
		}
	}
}

func (instance *PlanExecutor) executeJobs(jobs []Job) {
	var wg sync.WaitGroup
	for i := 0; i < instance.Config.Workers; i++ {
		wg.Add(1)
		go func(jobsForWorker []Job) {
			instance.workerExecuteJobs(jobsForWorker)
			wg.Done()
		}(jobs)
	}
	wg.Wait()
}

// Execute ...
func (instance *PlanExecutor) Execute() {
	instance.start = time.Now()
	instance.publisher = telegraph.NewLinkedPublisher()
	plan := instance.createPlan()

	go func() {
		instance.publisher.Publish(ExecutionStarted{})
		instance.executeJobs(plan.Jobs)
		instance.publisher.Publish(ExecutionStopped{})
	}()

	subscription := instance.publisher.Subscribe()
	for message := range subscription.Channel {
		switch message := message.(type) {
		case ExecutionStarted:
		case ExecutionStopped:
			subscription.RemoveFrom(instance.publisher)
		case ExecutionResult:
			executionResult := message
			for key, value := range executionResult {
				if handler, ok := resultHandlers[key]; ok {
					handler(value, instance.Stats)
				} else {
					logger.Log.Println(fmt.Sprintf("No handler for %s", key))
				}
			}
		default:
		}
	}
}
