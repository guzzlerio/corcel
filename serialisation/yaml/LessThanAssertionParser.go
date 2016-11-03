package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
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
