package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//YamlExactAssertionParser ...
type YamlExactAssertionParser struct{}

//Parse ...
func (instance YamlExactAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.ExactAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance YamlExactAssertionParser) Key() string {
	return "ExactAssertion"
}
