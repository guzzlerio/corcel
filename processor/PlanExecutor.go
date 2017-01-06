package processor

import (
	"context"
	"io/ioutil"
	"sync"
	"time"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/request"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"

	"github.com/REAANDREW/telegraph"
)

func merge(source map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
	for k, v := range extra {
		source[k] = v
	}
	return source
}

//PlanExecutionContext encapsulates the runtime state in order to execute
//a plan
type PlanExecutionContext struct {
	Plan         core.Plan
	Lists        *ListRingRevolver
	Config       *config.Configuration
	Publisher    telegraph.LinkedPublisher
	PlanContext  core.ExtractionResult
	JobContexts  map[int]core.ExtractionResult
	StepContexts map[int]map[int]core.ExtractionResult
	Bar          ProgressBar
	mutex        *sync.Mutex
	start        time.Time
}

//func (instance *PlanExecutionContext) execute(cancellation chan struct{}) {
func (instance *PlanExecutionContext) execute(ctx context.Context) {
	var jobs = instance.Plan.Jobs
	var jobStream = CreateJobStream(jobs, instance.Config)

	/*
		if int64(time.Since(instance.start).Seconds())%2 == 0 {
			instance.progress <- jobStream.Progress()
		}
	*/

	for jobStream.HasNext() {
		_ = instance.Bar.Set(jobStream.Progress())
		select {
		case <-ctx.Done():
			return
		default:
			job := jobStream.Next()

			for _, action := range job.Before {
				_ = action.Execute(ctx, nil)

			}
			instance.workerExecuteJob(ctx, job)
			for _, action := range job.After {
				_ = action.Execute(ctx, nil)
			}
		}
	}
}

func (instance *PlanExecutionContext) workerExecuteJob(ctx context.Context, job core.Job) {
	/*
		defer func() { //catch or finally
			if err := recover(); err != nil { //catch
				errormanager.Log(err)
			}
		}()
	*/
	var stepStream StepStream
	stepStream = CreateStepSequentialStream(job.Steps)
	if instance.Config.WaitTime > time.Duration(0) {
		stepStream = CreateStepDelayStream(stepStream, instance.Config.WaitTime)
	}

	for stepStream.HasNext() {
		step := stepStream.Next()
		//before Step
		for _, action := range step.Before {
			_ = action.Execute(ctx, nil)
		}
		executionResult := instance.executeStep(ctx, step)
		//instance.executeStep(ctx, step)

		//after Step
		for _, action := range step.After {
			_ = action.Execute(ctx, nil)
		}
		instance.Publisher.Publish(executionResult)
	}
}

func (instance *PlanExecutionContext) executeStep(ctx context.Context, step core.Step) core.ExecutionResult {
	start := time.Now()

	if instance.JobContexts[step.JobID] == nil {
		instance.JobContexts[step.JobID] = map[string]interface{}{}
		instance.StepContexts[step.JobID] = map[int]core.ExtractionResult{}
		instance.StepContexts[step.JobID][step.ID] = map[string]interface{}{}
	}

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
			vars := jValue.(map[string]interface{})
			for varKey, varValue := range vars {
				executionContext["$"+varKey] = varValue
			}
		}
	}

	var executionResult = core.ExecutionResult{}

	if step.Action != nil {
		executionResult = step.Action.Execute(ctx, executionContext)
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
	executionResult[core.DurationUrn.String()] = duration
	assertionResults := []core.AssertionResult{}
	for _, assertion := range step.Assertions {
		assertionResult := assertion.Assert(executionResult)
		assertionResults = append(assertionResults, assertionResult)
	}
	executionResult[core.AssertionsUrn.String()] = assertionResults

	return executionResult
}

//ExecutionBranch ...
type ExecutionBranch interface {
	Execute() error
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
	registry     core.Registry
	aggregator   statistics.AggregatorInterfaceToRenameLater
}

//CreatePlanExecutor ...
func CreatePlanExecutor(config *config.Configuration, bar ProgressBar, registry core.Registry, aggregator statistics.AggregatorInterfaceToRenameLater) *PlanExecutor {
	return &PlanExecutor{
		Config:       config,
		Bar:          bar,
		Publisher:    telegraph.NewLinkedPublisher(),
		PlanContext:  core.ExtractionResult{},
		JobContexts:  map[int]core.ExtractionResult{},
		StepContexts: map[int]map[int]core.ExtractionResult{},
		mutex:        &sync.Mutex{},
		registry:     registry,
		aggregator:   aggregator,
	}
}

//CreatePlanFromURLList ...
func CreatePlanFromURLList(config *config.Configuration) core.Plan {
	//FIXME Exposed for use in tests

	/*
		plan := core.Plan{
			Name:     "Plan from urls in file",
			Workers:  config.Workers,
			WaitTime: config.WaitTime,
			Jobs:     []core.Job{},
		}
	*/

	var name = "Plan from urls in file"
	var plan = core.NewPlanBuilder().
		Name(name).
		Workers(config.Workers).
		WaitTime(config.WaitTime).
		Build()

	reader := request.NewRequestReader(config.FilePath)

	stream := request.NewSequentialRequestStream(reader)

	for stream.HasNext() {
		job := plan.CreateJob()
		job.Name = "Job for the urls in file"

		request, err := stream.Next()
		if err != nil {
			errormanager.Check(err)
		}
		step := job.CreateStep()

		var body string
		if request.Body != nil {
			data, _ := ioutil.ReadAll(request.Body)
			if err != nil {
				errormanager.Check(err)
			} else {
				body = string(data)
			}
		}
		action := http.CreateAction()

		action.URL = request.URL.String()
		action.Method = request.Method
		action.Headers = request.Header
		action.Body = body

		step.Action = action
		//job.Steps = append(job.Steps, step)
		job = job.AddStep(step)
		//plan.Jobs = append(plan.Jobs, job)
		plan = plan.AddJob(job)
	}

	return plan
}

func (instance *PlanExecutor) generatePlan() core.Plan {
	var plan core.Plan
	var err error
	var config = instance.Config
	var registry = instance.registry

	if !config.Plan {
		plan = CreatePlanFromURLList(config)
	} else {
		parser := yaml.CreateExecutionPlanParser(registry)
		data, dataErr := ioutil.ReadFile(config.FilePath)
		if dataErr != nil {
			panic(dataErr)
		}
		plan, err = parser.Parse(string(data))
		config.Workers = plan.Workers
		if config.WaitTime == time.Duration(0) {
			config.WaitTime = plan.WaitTime
		}

		if config.Duration == time.Duration(0) {
			config.Duration = plan.Duration
		}

		if config.Iterations == 0 {
			config.Iterations = plan.Iterations
		}

		config.Random = plan.Random
		if err != nil {
			panic(err)
		}
	}
	return plan
}

// Execute ...
func (instance *PlanExecutor) Execute() error {
	//var cancellation = make(chan struct{})

	var ctx, cancel = context.WithCancel(context.Background())

	instance.start = time.Now()

	var mainPlan = instance.generatePlan()
	//before Plan
	for _, action := range mainPlan.Before {
		//_ = action.Execute(nil, cancellation)
		_ = action.Execute(ctx, nil)
	}

	var wg sync.WaitGroup
	wg.Add(instance.Config.Workers)

	var latch sync.WaitGroup
	latch.Add(instance.Config.Workers)

	var startingFence = make(chan struct{})

	var planChannel = make(chan core.Plan)

	for i := 0; i < instance.Config.Workers; i++ {
		go func() {
			defer errormanager.HandlePanic()
			//var plan = instance.generatePlan()
			var plan = <-planChannel
			if plan.Context["vars"] != nil {
				stringKeyData := map[string]interface{}{}
				data := plan.Context["vars"].(map[string]interface{})
				for dataKey, dataValue := range data {
					stringKeyData["$"+dataKey] = dataValue
				}
				plan.Context["vars"] = stringKeyData
			}

			var planExecutionContext = &PlanExecutionContext{
				Plan:         plan,
				Lists:        NewListRingRevolver(plan.Lists()),
				Config:       instance.Config,
				Publisher:    instance.Publisher,
				PlanContext:  core.ExtractionResult{},
				JobContexts:  map[int]core.ExtractionResult{},
				StepContexts: map[int]map[int]core.ExtractionResult{},
				Bar:          instance.Bar,
				mutex:        &sync.Mutex{},
				start:        time.Now(),
			}

			//planExecutionContext.execute(cancellation)
			latch.Done()
			<-startingFence
			planExecutionContext.execute(ctx)
			wg.Done()
		}()
		planChannel <- mainPlan
	}

	latch.Wait()

	instance.aggregator.Start()
	close(startingFence)
	if instance.Config.Duration > time.Duration(0) {
		time.AfterFunc(instance.Config.Duration, func() {
			//close(cancellation)
			cancel()
		})
	}
	wg.Wait()

	for _, action := range mainPlan.After {
		_ = action.Execute(ctx, nil)
	}

	return nil
}
