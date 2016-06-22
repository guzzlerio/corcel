package assertions

import (
	"fmt"
	"strconv"

	"ci.guzzler.io/guzzler/corcel/core"
)

//LessThanAssertion ...
type LessThanAssertion struct {
	Key   string
	Value interface{}
}

func (instance *LessThanAssertion) resultKey() string {
	return fmt.Sprintf("assert:lt:%v:%v", instance.Key, instance.Value)
}

//Assert ...
func (instance *LessThanAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"expected": instance.Value,
		"actual":   actual,
		"key":      instance.resultKey(),
	}

	// INEFFICIENT BUT WORKING ...
	//var actualType interface{}
	//var instanceType interface{}
	switch actualType := actual.(type) {
	case float64:
		switch instanceType := instance.Value.(type) {
		case float64:
			result["result"] = actualType < instanceType
			break
		case int:
			result["result"] = actualType < float64(instanceType)
			break
		case string:
			value, err := strconv.ParseFloat(instanceType, 64)
			if err != nil {
				result["result"] = false
			} else {
				result["result"] = actualType < value
			}
		default:
			result["result"] = false
		}
	case int:
		switch instanceType := instance.Value.(type) {
		case float64:
			result["result"] = float64(actualType) < instanceType
		case int:
			result["result"] = actualType < instanceType
		case string:
			value, err := strconv.ParseFloat(instanceType, 64)
			if err != nil {
				result["result"] = false
			} else {
				result["result"] = float64(actualType) < value
			}
		default:
			result["result"] = false
		}
	case string:
		switch instanceType := instance.Value.(type) {
		case float64:
			value, err := strconv.ParseFloat(actualType, 64)
			if err != nil {
				result["result"] = false
			} else {
				result["result"] = value < instanceType
			}
		case int:
			value, err := strconv.ParseFloat(actualType, 64)
			if err != nil {
				result["result"] = false
			} else {
				result["result"] = value < float64(instanceType)
			}
		case string:
			actualStrValue, actualStrErr := strconv.ParseFloat(actualType, 64)
			instanceStrValue, instanceStrErr := strconv.ParseFloat(instanceType, 64)
			if actualStrErr != nil && instanceStrErr != nil {
				result["result"] = actualType < instanceType
			} else {
				result["result"] = actualStrValue < instanceStrValue
			}
		default:
			result["result"] = false
		}
	default:
		result["result"] = false

	}

	if result["result"].(bool) == false {
		result["message"] = fmt.Sprintf("FAIL: %v is not greater %v", actual, instance.Value)
	}
	return result
}
