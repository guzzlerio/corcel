package yaml

import (
	"ci.guzzler.io/guzzler/corcel/assertions"
	"ci.guzzler.io/guzzler/corcel/core"
)

//LessThanOrEqualAssertionParser ...
type LessThanOrEqualAssertionParser struct{}

//Parse ...
func (instance LessThanOrEqualAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.LessThanOrEqualAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance LessThanOrEqualAssertionParser) Key() string {
	return "LessThanOrEqualAssertion"
}
