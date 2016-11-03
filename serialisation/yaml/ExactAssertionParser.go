package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
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
