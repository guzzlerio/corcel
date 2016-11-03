package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//GreaterThanOrEqualAssertionParser ...
type GreaterThanOrEqualAssertionParser struct{}

//Parse ...
func (instance GreaterThanOrEqualAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.GreaterThanOrEqualAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance GreaterThanOrEqualAssertionParser) Key() string {
	return "GreaterThanOrEqualAssertion"
}
