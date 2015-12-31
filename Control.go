package main

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/processor"
)

// Control ...
type Control interface {
	Start(*config.Configuration) (*ExecutionID, error)
	Stop(*ExecutionID) processor.ExecutionOutput
	Status(*ExecutionID) processor.ExecutionOutput
	History() []*ExecutionID
	Events() <-chan string

	//TODO Shouldn't need to expose this out, but required for transition
	Statistics() processor.Statistics
}

// Controller ...
type Controller struct {
	stats      *processor.Statistics
	executions map[*ExecutionID]*Executor
	bar        ProgressBar
}

// Start ...
func (instance *Controller) Start(config *config.Configuration) (*ExecutionID, error) {
	id := NewExecutionID()
	fmt.Printf("Execution ID: %s\n", id)
	executor := Executor{config, instance.stats, instance.bar}
	instance.executions[&id] = &executor
	executor.Execute()
	return &id, nil
}

// Stop ...
func (instance *Controller) Stop(id *ExecutionID) processor.ExecutionOutput {
	return instance.executions[id].Output()
}

// Status ...
func (instance *Controller) Status(*ExecutionID) processor.ExecutionOutput {
	return processor.ExecutionOutput{}
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
func (instance *Controller) Statistics() processor.Statistics {
	return *instance.stats
}

// NewControl ...
func NewControl(bar ProgressBar) Control {
	stats := processor.CreateStatistics()
	control := Controller{stats: stats, executions: make(map[*ExecutionID]*Executor), bar: bar}
	return &control
}
