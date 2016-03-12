package processor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/rcrowley/go-metrics"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/statistics"
)

// Control ...
type Control interface {
	Start(*config.Configuration) (*ExecutionID, error)
	Stop(*ExecutionID) ExecutionOutput
	Status(*ExecutionID) ExecutionOutput
	History() []*ExecutionID
	Events() <-chan string

	//TODO Shouldn't need to expose this out, but required for transition
	Statistics() Statistics
}

// Controller ...
type Controller struct {
	stats      *Statistics
	executions map[*ExecutionID]ExecutionBranch
	bar        ProgressBar
}

func (instance *Controller) createExecutionBranch(config *config.Configuration) ExecutionBranch {
	useNew := true

	if useNew {
		return CreatePlanExecutor(config, instance.stats, instance.bar)
	}

	return &Executor{config, instance.stats, instance.bar}
}

// Start ...
/*
func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()
	fmt.Printf("Execution ID: %s\n", id)

	executor := instance.createExecutionBranch(config)

	instance.executions[&id] = executor
	err := executor.Execute()
	return &id, err
}
*/

func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()
	fmt.Printf("Execution ID: %s\n", id)
	resultProcessors := []ExecutionResultProcessor{
		NewHTTPExecutionResultProcessor(),
		NewGeneralExecutionResultProcessor(),
	}

	aggregator := statistics.NewAggregator(metrics.DefaultRegistry)

	executor := CreatePlanExecutor(config, instance.stats, instance.bar)

	subscription := executor.Publisher.Subscribe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for executionResult := range subscription.Channel {
			result := executionResult.(ExecutionResult)
			for _, processor := range resultProcessors {
				processor.Process(result, metrics.DefaultRegistry, instance.stats)
			}
		}
		wg.Done()
	}()
	instance.executions[&id] = executor
	aggregator.Start()
	err := executor.Execute()
	subscription.RemoveFrom(executor.Publisher)
	wg.Wait()
	aggregator.Stop()

	b, err := json.Marshal(aggregator.Snapshot())
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./output.json", b, 0644)
	if err != nil {
		panic(err)
	}
	return &id, err
}

// Stop ...
func (instance *Controller) Stop(id *ExecutionID) ExecutionOutput {
	return instance.executions[id].Output()
}

// Status ...
func (instance *Controller) Status(*ExecutionID) ExecutionOutput {
	return ExecutionOutput{}
}

// History ...
func (instance *Controller) History() []*ExecutionID {
	return nil
}

// Events ...
func (instance *Controller) Events() <-chan string {
	return nil
}

// Statistics ...
func (instance *Controller) Statistics() Statistics {
	return *instance.stats
}

// NewControl ...
func NewControl(bar ProgressBar) Control {
	stats := CreateStatistics()
	executions := make(map[*ExecutionID]ExecutionBranch)
	control := Controller{stats: stats, executions: executions, bar: bar}
	return &control
}
