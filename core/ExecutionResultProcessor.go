package core

import "github.com/rcrowley/go-metrics"

//ExecutionResultProcessor ...
type ExecutionResultProcessor interface {
	Process(result ExecutionResult, registry metrics.Registry)
}
