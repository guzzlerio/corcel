package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//NotEqualAssertionParser ...
type NotEqualAssertionParser struct{}

//Parse ...
func (instance NotEqualAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.NotEqualAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance NotEqualAssertionParser) Key() string {
	return "NotEqualAssertion"
}
