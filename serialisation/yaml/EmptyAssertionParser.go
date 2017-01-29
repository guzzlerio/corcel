package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//EmptyAssertionParser ...
type EmptyAssertionParser struct{}

//Parse ...
func (instance EmptyAssertionParser) Parse(input map[string]interface{}) (core.Assertion, error) {

	if _, ok := input["key"]; !ok {
		return nil, errors.New("key is not present")
	}

	return &assertions.EmptyAssertion{
		Key: input["key"].(string),
	}, nil
}

//Key ...
func (instance EmptyAssertionParser) Key() string {
	return "EmptyAssertion"
}
