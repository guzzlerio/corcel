package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//EmptyAssertionParser ...
type EmptyAssertionParser struct{}

//Parse ...
func (instance EmptyAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.EmptyAssertion{
		Key: input["key"].(string),
	}
}

//Key ...
func (instance EmptyAssertionParser) Key() string {
	return "EmptyAssertion"
}
