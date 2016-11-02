package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//ExactAssertionParser ...
type ExactAssertionParser struct{}

//Parse ...
func (instance ExactAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.ExactAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance ExactAssertionParser) Key() string {
	return "ExactAssertion"
}
