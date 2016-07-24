package processor

import (
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/errormanager"

	"github.com/REAANDREW/telegraph"
)

func merge(source map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
	for k, v := range extra {
		source[k] = v
	}
	return source
}

//ExecutionBranch ...
type ExecutionBranch interface {
	Execute(plan core.Plan) error
}

//PlanExecutor ...
type PlanExecutor struct {
	Config       *config.Configuration
	Bar          ProgressBar
	start        time.Time
	Publisher    telegraph.LinkedPublisher
	Plan         core.Plan
	PlanContext  core.ExtractionResult
	JobContexts  map[int]core.ExtractionResult
	StepContexts map[int]map[int]core.ExtractionResult
	mutex        *sync.Mutex
}

//CreatePlanExecutor ...
func CreatePlanExecutor(config *config.Configuration, bar ProgressBar) *PlanExecutor {
	return &PlanExecutor{
		Config:       config,
		Bar:          bar,
		Publisher:    telegraph.NewLinkedPublisher(),
		PlanContext:  core.ExtractionResult{},
		JobContexts:  map[int]core.ExtractionResult{},
		StepContexts: map[int]map[int]core.ExtractionResult{},
		mutex:        &sync.Mutex{},
	}
}

func (instance *PlanExecutor) executeStep(step core.Step, cancellation chan struct{}) core.ExecutionResult {
	start := time.Now()
	instance.mutex.Lock()
	if instance.JobContexts[step.JobID] == nil {
		instance.JobContexts[step.JobID] = map[string]interface{}{}
		instance.StepContexts[step.JobID] = map[int]core.ExtractionResult{}
		instance.StepContexts[step.JobID][step.ID] = map[string]interface{}{}
	}
	instance.mutex.Unlock()

	var executionContext = core.ExecutionContext{}

	for pKey, pValue := range instance.Plan.Context {
		executionContext[pKey] = pValue
	}

	job := instance.Plan.GetJob(step.JobID)
	for jKey, jValue := range job.Context {
		executionContext[jKey] = jValue
	}

	var executionResult = core.ExecutionResult{}

	if step.Action != nil {
		executionResult = step.Action.Execute(executionContext, cancellation)
	}

	executionResult = merge(executionResult, instance.PlanContext)
	executionResult = merge(executionResult, instance.JobContexts[step.JobID])
	executionResult = merge(executionResult, instance.StepContexts[step.JobID][step.ID])

	executionResult = merge(executionResult, executionContext)

	for _, extractor := range step.Extractors {
		extractorResult := extractor.Extract(executionResult)

		switch extractorResult.Scope() {
		case core.PlanScope:
			instance.PlanContext = merge(instance.PlanContext, extractorResult)
			fallthrough
		case core.JobScope:
			instance.JobContexts[step.JobID] = merge(instance.JobContexts[step.JobID], extractorResult)
			fallthrough
		case core.StepScope:
			instance.StepContexts[step.JobID][step.ID] = merge(instance.StepContexts[step.JobID][step.ID], extractorResult)
		}

		//instance.JobContexts[step.JobID] = merge(instance.JobContexts[step.JobID], extractorResult)
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

	if instance.Config.Iterations > 0 {
		revolvingStream := CreateJobRevolvingStream(jobStream)
		iterationStream := CreateJobIterationStream(*revolvingStream, len(jobs), instance.Config.Iterations)
		jobStream = iterationStream
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
	instance.Plan = plan
	instance.executeJobs(plan)

	return nil
}
