package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
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
var resultHandlers = map[string]func(obj interface{}, statistics *Statistics){
	"http:request:error": func(obj interface{}, statistics *Statistics) {
		statistics.Request(obj.(error))
	},
	"http:response:error": func(obj interface{}, statistics *Statistics) {
		statistics.Request(obj.(error))
	},
	"http:request:bytes": func(obj interface{}, statistics *Statistics) {
		statistics.BytesSent(int64(obj.(int)))
		histogram := metrics.GetOrRegisterHistogram("http:request:bytes", metrics.DefaultRegistry, metrics.NewUniformSample(100))
		histogram.Update(int64(obj.(int)))
	},
	"http:response:bytes": func(obj interface{}, statistics *Statistics) {
		statistics.BytesReceived(int64(obj.(int)))
		histogram := metrics.GetOrRegisterHistogram("http:response:bytes", metrics.DefaultRegistry, metrics.NewUniformSample(100))
		histogram.Update(int64(obj.(int)))
	},
	"http:response:status": func(obj interface{}, statistics *Statistics) {
		statusCode := obj.(int)
		counter := metrics.GetOrRegisterCounter(fmt.Sprintf("http:response:status:%d", statusCode), metrics.DefaultRegistry)
		counter.Inc(1)

		var responseErr error
		if statusCode >= 400 && statusCode < 600 {
			responseErr = errors.New("5XX Response Code")
		}
		statistics.Request(responseErr)
	},
	"action:duration": func(obj interface{}, statistics *Statistics) {
		statistics.ResponseTime(int64(obj.(time.Duration)))
		//timer := metrics.GetOrRegisterTimer("action:duration", metrics.DefaultRegistry)
		//timer.Update(obj.(time.Duration))
	},
}

func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()
	fmt.Printf("Execution ID: %s\n", id)

	aggregator := statistics.NewAggregator(metrics.DefaultRegistry)

	executor := CreatePlanExecutor(config, instance.stats, instance.bar)

	subscription := executor.Publisher.Subscribe()
	count := 0
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for executionResult := range subscription.Channel {
			count = count + 1
			for key, value := range executionResult.(ExecutionResult) {
				if handler, ok := resultHandlers[key]; ok {
					handler(value, instance.stats)
				} else {
					logger.Log.Println(fmt.Sprintf("No handler for %s", key))
				}
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
