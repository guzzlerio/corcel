package core

//ExecutionActionParser ...
type ExecutionActionParser interface {
	Parse(input map[string]interface{}) Action
	Key() string
}
