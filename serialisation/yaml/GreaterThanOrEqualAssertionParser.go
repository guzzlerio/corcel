package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
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
