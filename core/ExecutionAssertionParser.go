package core

//ExecutionAssertionParser ...
type ExecutionAssertionParser interface {
	Parse(input map[string]interface{}) Assertion
	Key() string
}
