package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//NotEmptyAssertionParser ...
type NotEmptyAssertionParser struct{}

//Parse ...
func (instance NotEmptyAssertionParser) Parse(input map[string]interface{}) (core.Assertion, error) {

	if _, ok := input["key"]; !ok {
		return nil, errors.New("key is not present")
	}

	return &assertions.NotEmptyAssertion{
		Key: input["key"].(string),
	}, nil
}

//Key ...
func (instance NotEmptyAssertionParser) Key() string {
	return "NotEmptyAssertion"
}
