package assertions

import (
	"fmt"
	"strconv"

	. "github.com/smartystreets/goconvey/convey"
)

var NilValue interface{}

type AssertionTestCase struct {
	Actual               interface{}
	Instance             interface{}
	ActualStringNumber   bool
	InstanceStringNumber bool
}

func NewAsssertionTestCase(actual interface{}, instance interface{}) (newInstance AssertionTestCase) {
	newInstance.Actual = actual
	newInstance.Instance = instance
	switch actualType := actual.(type) {
	case string:
		_, err := strconv.ParseFloat(actualType, 64)
		if err == nil {
			newInstance.ActualStringNumber = true
		}
	}
	switch instanceType := instance.(type) {
	case string:
		_, err := strconv.ParseFloat(instanceType, 64)
		if err == nil {
			newInstance.InstanceStringNumber = true
		}
	}
	return
}

func assert(testCases []AssertionTestCase, test func(actual interface{}, instance interface{})) {

	for _, testCase := range testCases {
		actualValue := testCase.Actual
		instanceValue := testCase.Instance
		testName := fmt.Sprintf("When Actual is of type %T %[1]V and Instance is of type %T %[2]V", actualValue, instanceValue)
		if testCase.ActualStringNumber {
			testName = fmt.Sprintf("%s. Actual value is a STRING NUMBER in this case", testName)
		}
		if testCase.InstanceStringNumber {
			testName = fmt.Sprintf("%s. Instance value is a STRING NUMBER in this case", testName)
		}
		Convey(testName, func() {
			test(actualValue, instanceValue)
		})
	}
}
