package processor

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/request"
	"ci.guzzler.io/guzzler/corcel/statistics"
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
}

func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()
	resultProcessors := []ExecutionResultProcessor{
		NewHTTPExecutionResultProcessor(),
		NewGeneralExecutionResultProcessor(),
	}

	instance.aggregator = statistics.NewAggregator(metrics.DefaultRegistry)

	executor := CreatePlanExecutor(config, instance.bar)

	subscription := executor.Publisher.Subscribe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for executionResult := range subscription.Channel {
			result := executionResult.(ExecutionResult)
			for _, processor := range resultProcessors {
				processor.Process(result, metrics.DefaultRegistry)
			}
		}
		wg.Done()
	}()
	instance.executions[&id] = executor
	plan := GetPlan(config)
	instance.aggregator.Start()
	err := executor.Execute(plan)
	subscription.RemoveFrom(executor.Publisher)
	wg.Wait()
	return &id, err
}

func GetPlan(config *config.Configuration) Plan {
	var plan Plan
	var err error
	if !config.Plan {
		plan = CreatePlanFromConfiguration(config)
	} else {
		parser := CreateExecutionPlanParser()
		data, dataErr := ioutil.ReadFile(config.FilePath)
		if dataErr != nil {
			panic(dataErr)
		}
		plan, err = parser.Parse(string(data))
		config.Workers = plan.Workers
		if config.WaitTime == time.Duration(0) {
			config.WaitTime = plan.WaitTime
		}
		if err != nil {
			panic(err)
		}
	}
	return plan
}

func CreatePlanFromConfiguration(config *config.Configuration) Plan {
	plan := Plan{
		Name:     "Plan from urls in file",
		Workers:  config.Workers,
		WaitTime: config.WaitTime,
		Jobs:     []Job{},
	}

	reader := request.NewRequestReader(config.FilePath)

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
func NewControl(bar ProgressBar) Control {
	executions := make(map[*ExecutionID]ExecutionBranch)
	control := Controller{executions: executions, bar: bar}
	return &control
}
