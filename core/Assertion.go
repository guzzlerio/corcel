package core

//Assertion ...
type Assertion interface {
	Assert(ExecutionResult) AssertionResult
}
