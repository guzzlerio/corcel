package processor

import (
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/errormanager"

	"github.com/REAANDREW/telegraph"
)

//ExecutionBranch ...
type ExecutionBranch interface {
	Execute(plan core.Plan) error
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

func (instance *PlanExecutor) executeStep(step core.Step, cancellation chan struct{}) core.ExecutionResult {
	start := time.Now()
	executionResult := step.Action.Execute(cancellation)

	for _, extractor := range step.Extractors {
		extractorResult := extractor.Extract(executionResult)
		for k, v := range extractorResult {
			executionResult[k] = v
		}
	}

	duration := time.Since(start) / time.Millisecond
	executionResult["action:duration"] = duration
	assertionResults := []core.AssertionResult{}
	for _, assertion := range step.Assertions {
		assertionResult := assertion.Assert(executionResult)
		assertionResults = append(assertionResults, assertionResult)
	}
	executionResult["assertions"] = assertionResults

	return executionResult
}

func (instance *PlanExecutor) workerExecuteJob(talula core.Job, cancellation chan struct{}) {
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
		executionResult := instance.executeStep(step, cancellation)

		instance.Publisher.Publish(executionResult)

	}
}

func (instance *PlanExecutor) workerExecuteJobs(jobs []core.Job) {
	var jobStream JobStream
	jobStream = CreateJobSequentialStream(jobs)

	var cancellation = make(chan struct{})

	if instance.Config.Random {
		jobStream = CreateJobRandomStream(jobs)
	}
	if instance.Config.Duration > time.Duration(0) {
		jobStream = CreateJobDurationStream(jobStream, instance.Config.Duration)
		ticker := time.NewTicker(time.Millisecond * 10)
		go func() {
			for _ = range ticker.C {
				_ = instance.Bar.Set(jobStream.Progress())
			}
		}()
		time.AfterFunc(instance.Config.Duration, func() {
			ticker.Stop()
			_ = instance.Bar.Set(100)
			close(cancellation)
		})
	}

	for jobStream.HasNext() {
		job := jobStream.Next()
		_ = instance.Bar.Set(jobStream.Progress())
		instance.workerExecuteJob(job, cancellation)
	}
}

func (instance *PlanExecutor) executeJobs(plan core.Plan) {
	var wg sync.WaitGroup
	wg.Add(instance.Config.Workers)
	for i := 0; i < instance.Config.Workers; i++ {
		go func(executionPlan core.Plan) {
			instance.workerExecuteJobs(executionPlan.Jobs)
			wg.Done()
		}(plan)
	}
	wg.Wait()
}

// Execute ...
func (instance *PlanExecutor) Execute(plan core.Plan) error {
	instance.start = time.Now()
	instance.executeJobs(plan)

	return nil
}
