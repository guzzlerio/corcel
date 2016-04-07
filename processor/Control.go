package processor

import (
	"fmt"
	"sync"

	"github.com/rcrowley/go-metrics"

	"ci.guzzler.io/guzzler/corcel/config"
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
	fmt.Printf("Execution ID: %s\n", id)
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
	instance.aggregator.Start()
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
func NewControl(bar ProgressBar) Control {
	executions := make(map[*ExecutionID]ExecutionBranch)
	control := Controller{executions: executions, bar: bar}
	return &control
}
