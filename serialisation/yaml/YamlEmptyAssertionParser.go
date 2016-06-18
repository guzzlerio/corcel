package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//YamlEmptyAssertionParser ...
type YamlEmptyAssertionParser struct{}

//Parse ...
func (instance YamlEmptyAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.EmptyAssertion{
		Key: input["key"].(string),
	}
}

//Key ...
func (instance YamlEmptyAssertionParser) Key() string {
	return "EmptyAssertion"
}
