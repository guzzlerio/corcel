package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//GreaterThanOrEqualAssertionParser ...
type GreaterThanOrEqualAssertionParser struct{}

//Parse ...
func (instance GreaterThanOrEqualAssertionParser) Parse(input map[string]interface{}) (core.Assertion, error) {

	if _, ok := input["key"]; !ok {
		return nil, errors.New("key is not present")
	}

	if _, ok := input["expected"]; !ok {
		return nil, errors.New("expected is not present")
	}

	return &assertions.GreaterThanOrEqualAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}
}

//Key ...
func (instance GreaterThanOrEqualAssertionParser) Key() string {
	return "GreaterThanOrEqualAssertion"
}
