package assertions

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/core"
)

//ExactAssertion ...
type ExactAssertion struct {
	Key      string
	Expected interface{}
}

//ResultKey ...
func (instance *ExactAssertion) ResultKey() string {
	return fmt.Sprintf("assert:exactmatch:%v:%v", instance.Key, instance.Expected)
}

//Assert ...
func (instance *ExactAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"expected": instance.Expected,
		"actual":   actual,
		"key":      instance.ResultKey(),
	}
	if actual == instance.Expected {
		result["result"] = true
	} else {
		result["result"] = false
		result["message"] = fmt.Sprintf("FAIL: %v does not match %v", actual, instance.Expected)
	}
	return result
}
