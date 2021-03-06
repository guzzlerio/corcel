package assertions

import (
	"fmt"
	"strconv"

	"github.com/guzzlerio/corcel/core"
)

//GreaterThanAssertion ...
type GreaterThanAssertion struct {
	Key   string
	Value interface{}
}

func (instance *GreaterThanAssertion) resultKey() string {
	return fmt.Sprintf("assert:gt:%v:%v", instance.Key, instance.Value)
}

//Assert ...
func (instance *GreaterThanAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"expected": instance.Value,
		"actual":   actual,
		"key":      instance.resultKey(),
	}

	switch actualType := actual.(type) {
	case float64:
		switch instanceType := instance.Value.(type) {
		case float64:
			result[core.AssertionResultUrn.String()] = actualType > instanceType
			break
		case int:
			result[core.AssertionResultUrn.String()] = actualType > float64(instanceType)
			break
		case string:
			value, err := strconv.ParseFloat(instanceType, 64)
			if err != nil {
				result[core.AssertionResultUrn.String()] = false
			} else {
				result[core.AssertionResultUrn.String()] = actualType > value
			}
		default:
			result[core.AssertionResultUrn.String()] = true
		}
	case int:
		switch instanceType := instance.Value.(type) {
		case float64:
			result[core.AssertionResultUrn.String()] = float64(actualType) > instanceType
		case int:
			result[core.AssertionResultUrn.String()] = actualType > instanceType
		case string:
			value, err := strconv.ParseFloat(instanceType, 64)
			if err != nil {
				result[core.AssertionResultUrn.String()] = false
			} else {
				result[core.AssertionResultUrn.String()] = float64(actualType) > value
			}
		default:
			result[core.AssertionResultUrn.String()] = true
		}
	case string:
		switch instanceType := instance.Value.(type) {
		case float64:
			value, err := strconv.ParseFloat(actualType, 64)
			if err != nil {
				result[core.AssertionResultUrn.String()] = false
			} else {
				result[core.AssertionResultUrn.String()] = value > instanceType
			}
		case int:
			value, err := strconv.ParseFloat(actualType, 64)
			if err != nil {
				result[core.AssertionResultUrn.String()] = false
			} else {
				result[core.AssertionResultUrn.String()] = value > float64(instanceType)
			}
		case string:
			actualStrValue, actualStrErr := strconv.ParseFloat(actualType, 64)
			instanceStrValue, instanceStrErr := strconv.ParseFloat(instanceType, 64)
			if actualStrErr != nil && instanceStrErr != nil {
				result[core.AssertionResultUrn.String()] = actualType > instanceType
			} else {
				result[core.AssertionResultUrn.String()] = actualStrValue > instanceStrValue
			}
		default:
			result[core.AssertionResultUrn.String()] = true
		}
	default:
		result[core.AssertionResultUrn.String()] = false

	}

	if result[core.AssertionResultUrn.String()].(bool) == false {
		result[core.AssertionMessageUrn.String()] = fmt.Sprintf("FAIL: %v is not greater %v", actual, instance.Value)
	}
	return result
}
