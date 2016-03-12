package processor

import (
	"io/ioutil"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/request"

	"github.com/REAANDREW/telegraph"
)

//PlanExecutor ...
type PlanExecutor struct {
	Config    *config.Configuration
	Stats     *Statistics
	Bar       ProgressBar
	start     time.Time
	Publisher telegraph.LinkedPublisher
}

//CreatePlanExecutor ...
func CreatePlanExecutor(config *config.Configuration, stats *Statistics, bar ProgressBar) *PlanExecutor {
	return &PlanExecutor{
		Config:    config,
		Bar:       bar,
		Stats:     stats,
		Publisher: telegraph.NewLinkedPublisher(),
	}
}

func (instance *PlanExecutor) createPlan() Plan {
	plan := Plan{
		Name:     "Plan from urls in file",
		Workers:  instance.Config.Workers,
		WaitTime: instance.Config.WaitTime,
		Jobs:     []Job{},
	}

	reader := request.NewRequestReader(instance.Config.FilePath)

	stream := request.NewSequentialRequestStream(reader)

	for stream.HasNext() {
		job := Job{
			Name:  "Job for the urls in file",
			Steps: []Step{},
		}

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
		plan.Jobs = append(plan.Jobs, job)
	}

	return plan
}

/*
var resultHandlers = map[string]func(obj interface{}, statistics *Statistics){
	"http:request:error": func(obj interface{}, statistics *Statistics) {
		statistics.Request(obj.(error))
	},
	"http:response:error": func(obj interface{}, statistics *Statistics) {
		statistics.Request(obj.(error))
	},
	"http:request:bytes": func(obj interface{}, statistics *Statistics) {
		statistics.BytesSent(int64(obj.(int)))
		histogram := metrics.GetOrRegisterHistogram("http:request:bytes", metrics.DefaultRegistry, metrics.NewUniformSample(100))
		histogram.Update(int64(obj.(int)))
	},
	"http:response:bytes": func(obj interface{}, statistics *Statistics) {
		statistics.BytesReceived(int64(obj.(int)))
		histogram := metrics.GetOrRegisterHistogram("http:response:bytes", metrics.DefaultRegistry, metrics.NewUniformSample(100))
		histogram.Update(int64(obj.(int)))
	},
	"http:response:status": func(obj interface{}, statistics *Statistics) {
		statusCode := obj.(int)
		counter := metrics.GetOrRegisterCounter(fmt.Sprintf("http:response:status:%d", statusCode), metrics.DefaultRegistry)
		counter.Inc(1)

		var responseErr error
		if statusCode >= 400 && statusCode < 600 {
			responseErr = errors.New("5XX Response Code")
		}
		statistics.Request(responseErr)
	},
	"action:duration": func(obj interface{}, statistics *Statistics) {
		statistics.ResponseTime(int64(obj.(time.Duration)))
		//timer := metrics.GetOrRegisterTimer("action:duration", metrics.DefaultRegistry)
		//timer.Update(obj.(time.Duration))
	},
}
*/

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

func (instance *PlanExecutor) workerExecuteJob(talula Job) {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			errormanager.Log(err)
		}
	}()
	var stepStream StepStream
	stepStream = CreateStepSequentialStream(talula.Steps)
	if instance.Config.WaitTime > time.Duration(0) {
		stepStream = CreateStepDelayStream(stepStream, instance.Config.WaitTime)
	}

	for stepStream.HasNext() {
		step := stepStream.Next()
		executionResult := instance.executeStep(step)

		instance.Publisher.Publish(executionResult)

	}
}

func (instance *PlanExecutor) workerExecuteJobs(jobs []Job) {
	var jobStream JobStream
	jobStream = CreateJobSequentialStream(jobs)

	if instance.Config.Random {
		jobStream = CreateJobRandomStream(jobs)
	}
	if instance.Config.Duration > time.Duration(0) {
		jobStream = CreateJobDurationStream(jobStream, instance.Config.Duration)
	}

	for jobStream.HasNext() {
		job := jobStream.Next()
		_ = instance.Bar.Set(jobStream.Progress())
		instance.workerExecuteJob(job)
	}
}

func (instance *PlanExecutor) executeJobs() {
	var wg sync.WaitGroup
	for i := 0; i < instance.Config.Workers; i++ {
		wg.Add(1)
		go func() {
			var plan Plan
			var err error
			if !instance.Config.Plan {
				plan = instance.createPlan()
			} else {
				parser := CreateExecutionPlanParser()
				data, dataErr := ioutil.ReadFile(instance.Config.FilePath)
				if dataErr != nil {
					panic(dataErr)
				}
				plan, err = parser.Parse(string(data))
				if err != nil {
					panic(err)
				}
			}
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
	instance.executeJobs()
	instance.Stats.Stop()

	return nil
}

// Output ...
func (instance *PlanExecutor) Output() ExecutionOutput {
	return instance.Stats.ExecutionOutput()
}
