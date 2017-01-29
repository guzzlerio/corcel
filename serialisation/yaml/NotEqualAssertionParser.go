package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//NotEqualAssertionParser ...
type NotEqualAssertionParser struct{}

//Parse ...
func (instance NotEqualAssertionParser) Parse(input map[string]interface{}) (core.Assertion, error) {

	if _, ok := input["key"]; !ok {
		return nil, errors.New("key is not present")
	}

	if _, ok := input["expected"]; !ok {
		return nil, errors.New("expected is not present")
	}

	return &assertions.NotEqualAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}, nil
}

//Key ...
func (instance NotEqualAssertionParser) Key() string {
	return "NotEqualAssertion"
}
