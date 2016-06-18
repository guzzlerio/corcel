package inproc

import "ci.guzzler.io/guzzler/corcel/core"

//DummyAction ...
type DummyAction struct {
	Results map[string]interface{}
}

//Execute ...
func (instance DummyAction) Execute(cancellation chan struct{}) core.ExecutionResult {
	result := core.ExecutionResult{}

	for key, value := range instance.Results {
		result[key] = value
	}

	return result
}
