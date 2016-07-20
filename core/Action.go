package core

//Action ...
type Action interface {
	Execute(executionContext ExecutionContext, cancellation chan struct{}) ExecutionResult
}
