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
			result[core.AssertionResultUrn.String()] = true
		}
	default:
		if actual != nil {
			result[core.AssertionResultUrn.String()] = true
		}
	}

	if result[core.AssertionResultUrn.String()] != true {
		result[core.AssertionResultUrn.String()] = false
		result[core.AssertionMessageUrn.String()] = fmt.Sprintf("FAIL: value is empty")
	}
	return result
}
