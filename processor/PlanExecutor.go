package processor

import (
	"fmt"
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
	Lists        *ListRingRevolver
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

	var vars map[string]interface{}

	if instance.Plan.Context["vars"] != nil {
		vars = instance.Plan.Context["vars"].(map[string]interface{})
	}

	for pKey, pValue := range vars {
		executionContext[pKey] = pValue
	}

	listValues := instance.Lists.Values()
	for pKey, pValue := range listValues {
		executionContext[pKey] = pValue
	}

	job := instance.Plan.GetJob(step.JobID)
	for jKey, jValue := range job.Context {
		executionContext[jKey] = jValue
		if jKey == "vars" {
			vars := jValue.(map[interface{}]interface{})
			for varKey, varValue := range vars {
				executionContext["$"+varKey.(string)] = varValue
			}
		}
	}

	var executionResult = core.ExecutionResult{}

	if step.Action != nil {
		executionResult = step.Action.Execute(executionContext, cancellation)
	}

	executionResult = merge(executionResult, instance.PlanContext)
	executionResult = merge(executionResult, instance.JobContexts[step.JobID])
	executionResult = merge(executionResult, instance.StepContexts[step.JobID][step.ID])

	executionResult = merge(executionResult, executionContext)

	instance.mutex.Lock()
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
	instance.mutex.Unlock()

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
		//before Step
		executionResult := instance.executeStep(step, cancellation)
		//after Step

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
		var vars map[string]interface{}

		if instance.Plan.Context["vars"] != nil {
			vars = job.Context["vars"].(map[string]interface{})
		}

		var executionContext core.ExecutionContext
		for pKey, pValue := range vars {
			executionContext[pKey] = pValue
		}

		listValues := instance.Lists.Values()
		for pKey, pValue := range listValues {
			executionContext[pKey] = pValue
		}
		_ = instance.Bar.Set(jobStream.Progress())
		//before Job
		for _, action := range job.Before {
			fmt.Println("Executing Before Job")
			_ = action.Execute(executionContext, nil)
		}
		instance.workerExecuteJob(job, cancellation)
		//after Job
		for _, action := range job.After {
			fmt.Println("Executing After Job")
			_ = action.Execute(executionContext, nil)
		}
	}
}

func (instance *PlanExecutor) executeJobs(jobs []core.Job) {
	var wg sync.WaitGroup
	wg.Add(instance.Config.Workers)
	for i := 0; i < instance.Config.Workers; i++ {
		go func(executionJobs []core.Job) {
			instance.workerExecuteJobs(executionJobs)
			wg.Done()
		}(jobs)
	}
	wg.Wait()
}

// Execute ...
func (instance *PlanExecutor) Execute(plan core.Plan) error {
	instance.start = time.Now()
	// fmt.Printf("Zee Plan: %+v", plan)
	instance.Plan = plan
	if instance.Plan.Context["lists"] != nil {
		var lists = map[string][]map[string]interface{}{}

		listKeys := instance.Plan.Context["lists"].(map[interface{}]interface{})
		for listKey, listValue := range listKeys {
			lists[listKey.(string)] = []map[string]interface{}{}
			listValueItems := listValue.([]interface{})
			for _, listValueItem := range listValueItems {
				srcData := listValueItem.(map[interface{}]interface{})
				stringKeyData := map[string]interface{}{}
				for srcKey, srcValue := range srcData {
					stringKeyData[srcKey.(string)] = srcValue
				}
				lists[listKey.(string)] = append(lists[listKey.(string)], stringKeyData)
			}
		}

		instance.Lists = NewListRingRevolver(lists)
	} else {
		instance.Lists = NewListRingRevolver(map[string][]map[string]interface{}{})
	}
	if instance.Plan.Context["vars"] != nil {
		stringKeyData := map[string]interface{}{}
		data := instance.Plan.Context["vars"].(map[interface{}]interface{})
		for dataKey, dataValue := range data {
			stringKeyData["$"+dataKey.(string)] = dataValue
		}
		instance.Plan.Context["vars"] = stringKeyData
	}
	//before Plan
	//TODO this is duplicated from executeStep. Extract
	var executionContext = core.ExecutionContext{}
	for pKey, pValue := range instance.Plan.Context {
		executionContext[pKey] = pValue
	}
	for _, action := range plan.Before {
		_ = action.Execute(executionContext, nil)
	}
	instance.executeJobs(plan.Jobs)
	//after Plan
	for _, action := range plan.After {
		_ = action.Execute(executionContext, nil)
	}

	return nil
}
