package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//GreaterThanAssertionParser ...
type GreaterThanAssertionParser struct{}

//Parse ...
func (instance GreaterThanAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.GreaterThanAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance GreaterThanAssertionParser) Key() string {
	return "GreaterThanAssertion"
}
