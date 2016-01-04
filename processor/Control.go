package processor

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/config"
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
	executions map[*ExecutionID]*Executor
	bar        ProgressBar
}

// Start ...
func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()
	fmt.Printf("Execution ID: %s\n", id)
	executor := Executor{config, instance.stats, instance.bar}
	instance.executions[&id] = &executor
	err := executor.Execute();
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
	control := Controller{stats: stats, executions: make(map[*ExecutionID]*Executor), bar: bar}
	return &control
}
