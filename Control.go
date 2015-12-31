package main

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/config"
)

// Control ...
type Control interface {
	Start(*config.Configuration) (*ExecutionId, error)
	Stop(*ExecutionId) ExecutionOutput
	Status(*ExecutionId) ExecutionOutput
	History() []*ExecutionId
	Events() <-chan string

	//TODO Shouldn't need to expose this out, but required for transition
	Statistics() Statistics
}

// Controller ...
type Controller struct {
	stats      *Statistics
	executions map[*ExecutionId]*Executor
	bar        ProgressBar
}

// Start ...
func (instance *Controller) Start(config *config.Configuration) (*ExecutionId, error) {
	id := NewExecutionId()
	fmt.Printf("Execution ID: %s\n", id)
	executor := Executor{config, instance.stats, instance.bar}
	instance.executions[&id] = &executor
	executor.Execute()
	return &id, nil
}

// Stop ...
func (instance *Controller) Stop(id *ExecutionId) ExecutionOutput {
	return instance.executions[id].Output()
}

// Status ...
func (instance *Controller) Status(*ExecutionId) ExecutionOutput {
	return ExecutionOutput{}
}

// History ...
func (instance *Controller) History() []*ExecutionId {
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
	control := Controller{stats: stats, executions: make(map[*ExecutionId]*Executor), bar: bar}
	return &control
}
