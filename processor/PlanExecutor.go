package processor

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"
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

//CreatePlanExecutor ...
func CreatePlanExecutor(config *config.Configuration, stats *Statistics, bar ProgressBar) *PlanExecutor {
	return &PlanExecutor{
		Config: config,
		Bar:    bar,
		Stats:  stats,
	}
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

	for stream.HasNext() {
		request, err := stream.Next()
		if err != nil {
			errormanager.Check(err)
		}
		step := Step{}

		action := &HTTPRequestExecutionAction{
			//Client:  client,
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
		statusCode := obj.(int)
		var responseErr error
		if statusCode >= 400 && statusCode < 600 {
			responseErr = errors.New("5XX Response Code")
		}
		statistics.Request(responseErr)
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
			defer func() { //catch or finally
				if err := recover(); err != nil { //catch
					errormanager.Log(err)
				}
			}()
			var stepStream StepStream
			stepStream = CreateStepSequentialStream(talula.Steps)
			if instance.Config.Random {
				stepStream = CreateStepRandomStream(talula.Steps)
			}
			if instance.Config.WaitTime > time.Duration(0) {
				stepStream = CreateStepDelayStream(stepStream, instance.Config.WaitTime)
			}

			if instance.Config.Duration > time.Duration(0) {
				stepStream = CreateStepDurationStream(stepStream, instance.Config.Duration)
			}

			for stepStream.HasNext() {
				_ = instance.Bar.Set(stepStream.Progress())
				step := stepStream.Next()
				executionResult := instance.executeStep(step)

				for key, value := range executionResult {
					if handler, ok := resultHandlers[key]; ok {
						handler(value, instance.Stats)
					} else {
						logger.Log.Println(fmt.Sprintf("No handler for %s", key))
					}
				}

				//instance.publisher.Publish(executionResult)
			}
		}(job)
		/*
			if instance.Config.Duration > 0 && time.Since(instance.start) < instance.Config.Duration {
				instance.workerExecuteJobs(jobs)
			} else {
				break
			}
		*/
	}
}

func (instance *PlanExecutor) executeJobs() {
	var wg sync.WaitGroup
	for i := 0; i < instance.Config.Workers; i++ {
		wg.Add(1)
		go func() {
			plan := instance.createPlan()
			instance.workerExecuteJobs(plan.Jobs)
			wg.Done()
		}()
	}
	wg.Wait()
}

// Execute ...
func (instance *PlanExecutor) Execute() error {
	instance.start = time.Now()
	instance.Stats.Start()
	//instance.publisher = telegraph.NewLinkedPublisher()

	//instance.publisher.Publish(ExecutionStarted{})
	instance.executeJobs()
	//instance.publisher.Publish(ExecutionStopped{})

	instance.Stats.Stop()
	return nil
}

// Output ...
func (instance *PlanExecutor) Output() ExecutionOutput {
	return instance.Stats.ExecutionOutput()
}
