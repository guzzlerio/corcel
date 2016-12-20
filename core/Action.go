package core

import "context"

//Action ...
type Action interface {
	//Execute(executionContext ExecutionContext, cancellation chan struct{}) ExecutionResult
	Execute(ctx context.Context, executionContext ExecutionContext) ExecutionResult
}
