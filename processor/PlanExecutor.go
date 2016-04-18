package processor

import (
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"

	"github.com/REAANDREW/telegraph"
)

//ExecutionBranch ...
type ExecutionBranch interface {
	Execute(plan Plan) error
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

func (instance *PlanExecutor) executeJobs(plan Plan) {
	var wg sync.WaitGroup
	wg.Add(instance.Config.Workers)
	for i := 0; i < instance.Config.Workers; i++ {
		go func(executionPlan Plan) {
			instance.workerExecuteJobs(executionPlan.Jobs)
			wg.Done()
		}(plan)
	}
	wg.Wait()
}

// Execute ...
func (instance *PlanExecutor) Execute(plan Plan) error {
	instance.start = time.Now()
	instance.executeJobs(plan)

	return nil
}