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

//ExecutionBranch ...
type ExecutionBranch interface {
	Execute() error
}

//PlanExecutor ...
type PlanExecutor struct {
	Config    *config.Configuration
	Bar       ProgressBar
	start     time.Time
	Publisher telegraph.LinkedPublisher
}

//CreatePlanExecutor ...
func CreatePlanExecutor(config *config.Configuration, bar ProgressBar) *PlanExecutor {
	return &PlanExecutor{
		Config:    config,
		Bar:       bar,
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
	instance.executeJobs()

	return nil
}
