package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
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
