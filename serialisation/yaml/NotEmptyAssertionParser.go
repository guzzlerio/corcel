package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//NotEmptyAssertionParser ...
type NotEmptyAssertionParser struct{}

//Parse ...
func (instance NotEmptyAssertionParser) Parse(input map[string]interface{}) core.Assertion {
	return &assertions.NotEmptyAssertion{
		Key: input["key"].(string),
	}
}

//Key ...
func (instance NotEmptyAssertionParser) Key() string {
	return "NotEmptyAssertion"
}
