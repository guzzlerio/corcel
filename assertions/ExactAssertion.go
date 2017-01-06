package assertions

import (
	"fmt"

	"github.com/guzzlerio/corcel/core"
)

//ExactAssertion ...
type ExactAssertion struct {
	Key   string
	Value interface{}
}

func (instance *ExactAssertion) resultKey() string {
	return fmt.Sprintf("assert:exactmatch:%v:%v", instance.Key, instance.Value)
}

func (instance *ExactAssertion) baseResult(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]
	return map[string]interface{}{
		"expected": instance.Value,
		"actual":   actual,
		"key":      instance.resultKey(),
	}
}

func (instance *ExactAssertion) pass(executionResult core.ExecutionResult) core.AssertionResult {
	var result = instance.baseResult(executionResult)
	result[core.AssertionResultUrn.String()] = true
	return result
}

func (instance *ExactAssertion) fail(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]
	var result = instance.baseResult(executionResult)
	var message = fmt.Sprintf("FAIL: %v %[1]T does not match %v %[2]T", actual, instance.Value)

	result[core.AssertionResultUrn.String()] = false
	result[core.AssertionMessageUrn.String()] = message
	return result
}

//Assert ...
func (instance *ExactAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
	actual := executionResult[instance.Key]

	//TODO:  Check the efficiency of this!!  This is by far the easiest way to do this but it might hinder performance
	//       Without doing this you will need to check for actual and expected types, look at the greater than assertion
	//       for an example.   Since the yaml library change and the fact that everything ends up in json, the types it uses
	//       are string, bool and float64 but not int.
	if fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", instance.Value) {
		return instance.pass(executionResult)
	}
	return instance.fail(executionResult)
}
