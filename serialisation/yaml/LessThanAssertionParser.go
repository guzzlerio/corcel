package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/assertions"
	"github.com/guzzlerio/corcel/core"
)

//LessThanAssertionParser ...
type LessThanAssertionParser struct{}

//Parse ...
func (instance LessThanAssertionParser) Parse(input map[string]interface{}) (core.Assertion, error) {

	if _, ok := input["key"]; !ok {
		return nil, errors.New("key is not present")
	}

	if _, ok := input["expected"]; !ok {
		return nil, errors.New("expected is not present")
	}

	return &assertions.LessThanAssertion{
		Key:   input["key"].(string),
		Value: input["expected"],
	}, nil
}

//Key ...
func (instance LessThanAssertionParser) Key() string {
	return "LessThanAssertion"
}
