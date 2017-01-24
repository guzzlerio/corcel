package processor

import (
	"sync"

	"github.com/rcrowley/go-metrics"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
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
	aggregator statistics.AggregatorInterfaceToRenameLater
	registry   core.Registry
}

//Start ...
func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()

	var metricsRegistry = metrics.NewRegistry()
	instance.aggregator = statistics.NewAggregator(metricsRegistry)

	executor := CreatePlanExecutor(config, instance.bar, instance.registry, instance.aggregator)

	subscription := executor.Publisher.Subscribe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer errormanager.HandlePanic()
		for executionResult := range subscription.Channel {
			//inproc.ProcessEventsSubscribed = inproc.ProcessEventsSubscribed + 1
			result := executionResult.(core.ExecutionResult)
			for _, processor := range instance.registry.ResultProcessors {
				processor.Process(result, metricsRegistry)
			}
		}
		wg.Done()
	}()
	instance.executions[&id] = executor
	//instance.aggregator.Start()
	err := executor.Execute()
	subscription.RemoveFrom(executor.Publisher)
	wg.Wait()
	return &id, err
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
