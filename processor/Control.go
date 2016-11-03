package processor

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/request"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
)

// Control ...
type Control interface {
	Start(*config.Configuration) (*ExecutionID, error)
	Stop(*ExecutionID) statistics.AggregatorSnapShot
	Status(*ExecutionID) statistics.AggregatorSnapShot
	History() []*ExecutionID
	Events() <-chan string
}

// Controller ...
type Controller struct {
	executions map[*ExecutionID]ExecutionBranch
	bar        ProgressBar
	aggregator *statistics.Aggregator
	registry   core.Registry
}

//Start ...
func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()

	instance.aggregator = statistics.NewAggregator(metrics.DefaultRegistry)

	executor := CreatePlanExecutor(config, instance.bar)

	subscription := executor.Publisher.Subscribe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for executionResult := range subscription.Channel {
			result := executionResult.(core.ExecutionResult)
			for _, processor := range instance.registry.ResultProcessors {
				processor.Process(result, metrics.DefaultRegistry)
			}
		}
		wg.Done()
	}()
	instance.executions[&id] = executor
	plan := getPlan(config, instance.registry)
	instance.aggregator.Start()
	err := executor.Execute(plan)
	subscription.RemoveFrom(executor.Publisher)
	wg.Wait()
	return &id, err
}

func getPlan(config *config.Configuration, registry core.Registry) core.Plan {
	var plan core.Plan
	var err error
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

//CreatePlanFromURLList ...
func CreatePlanFromURLList(config *config.Configuration) core.Plan {
	//FIXME Exposed for use in tests
	plan := core.Plan{
		Name:     "Plan from urls in file",
		Workers:  config.Workers,
		WaitTime: config.WaitTime,
		Jobs:     []core.Job{},
	}

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

// Stop ...
//A1
func (instance *Controller) Stop(id *ExecutionID) statistics.AggregatorSnapShot {
	instance.aggregator.Stop()

	return instance.aggregator.Snapshot()
}

// Status ...
func (instance *Controller) Status(*ExecutionID) statistics.AggregatorSnapShot {
	return instance.aggregator.Snapshot()
}

// History ...
func (instance *Controller) History() []*ExecutionID {
	return nil
}

// Events ...
func (instance *Controller) Events() <-chan string {
	return nil
}

// NewControl ...
func NewControl(bar ProgressBar, registry core.Registry) Control {
	//FIXME Possible no tests over the ExecutionBranch
	executions := make(map[*ExecutionID]ExecutionBranch)
	control := Controller{
		executions: executions,
		bar:        bar,
		registry:   registry,
	}
	return &control
}
