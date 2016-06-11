package core

//ExecutionResult ...
type ExecutionResult map[string]interface{}

//AssertionResult ...
type AssertionResult map[string]interface{}

//Action ...
type Action interface {
	Execute(cancellation chan struct{}) ExecutionResult
}

//Assertion ...
type Assertion interface {
	Assert(ExecutionResult) AssertionResult
}
