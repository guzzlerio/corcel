package assertions

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/core"
)

//NotEqualAssertion ...
type NotEqualAssertion struct {
	Key   string
	Value interface{}
}

func (instance *NotEqualAssertion) resultKey() string {
	return fmt.Sprintf("assert:notequal:%v:%v", instance.Key, instance.Value)
}

//Assert ...
func (instance *NotEqualAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"expected": instance.Value,
		"actual":   actual,
		"key":      instance.resultKey(),
	}

	if actual != instance.Value {
		result["result"] = true
	} else {
		result["result"] = false
		result["message"] = fmt.Sprintf("FAIL: %v does match %v", actual, instance.Value)
	}
	return result
}
