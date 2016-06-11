package assertions

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/core"
)

//ExactAssertion ...
type ExactAssertion struct {
	Key   string
	Value interface{}
}

func (instance *ExactAssertion) resultKey() string {
	return fmt.Sprintf("assert:exactmatch:%v:%v", instance.Key, instance.Value)
}

//Assert ...
func (instance *ExactAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"expected": instance.Value,
		"actual":   actual,
		"key":      instance.resultKey(),
	}
	if actual == instance.Value {
		result["result"] = true
	} else {
		result["result"] = false
		result["message"] = fmt.Sprintf("FAIL: %v does not match %v", actual, instance.Value)
	}
	return result
}
