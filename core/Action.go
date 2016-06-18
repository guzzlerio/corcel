package core

//Action ...
type Action interface {
	Execute(cancellation chan struct{}) ExecutionResult
}
