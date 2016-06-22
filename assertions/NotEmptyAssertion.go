package assertions

import (
	"fmt"
	"strings"

	"ci.guzzler.io/guzzler/corcel/core"
)

//NotEmptyAssertion ...
type NotEmptyAssertion struct {
	Key string
}

func (instance *NotEmptyAssertion) resultKey() string {
	return fmt.Sprintf("assert:notempty:%v", instance.Key)
}

//Assert ...
func (instance *NotEmptyAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"actual": actual,
		"key":    instance.resultKey(),
	}

	switch actualValue := actual.(type) {
	case string:
		value := strings.TrimSpace(actualValue)
		if value != "" {
			result["result"] = true
		}
	default:
		if actual != nil {
			result["result"] = true
		}
	}

	if result["result"] != true {
		result["result"] = false
		result["message"] = fmt.Sprintf("FAIL: value is empty")
	}
	return result
}
