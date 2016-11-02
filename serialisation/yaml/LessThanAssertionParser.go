package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//LessThanAssertionParser ...
type LessThanAssertionParser struct{}

//Parse ...
func (instance LessThanAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.LessThanAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance LessThanAssertionParser) Key() string {
	return "LessThanAssertion"
}
