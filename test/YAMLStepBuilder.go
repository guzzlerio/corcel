package test

import "ci.guzzler.io/guzzler/corcel/serialisation/yaml"

//YamlStepBuilder ...
type YamlStepBuilder struct {
	Action     map[string]interface{}
	Assertions []map[string]interface{}
	Extractors []map[string]interface{}
	Before     []yaml.Action
	After      []yaml.Action
}

//Build ...
func (instance *YamlStepBuilder) Build() yaml.ExecutionStep {
	return yaml.ExecutionStep{
		Action:     instance.Action,
		Assertions: instance.Assertions,
		Extractors: instance.Extractors,
		Before:     instance.Before,
		After:      instance.After,
	}
}

//CreateStepBuilder ...
func CreateStepBuilder() *YamlStepBuilder {
	return &YamlStepBuilder{
		Action:     map[string]interface{}{},
		Assertions: []map[string]interface{}{},
		Before:     []yaml.Action{},
		After:      []yaml.Action{},
	}
}

//AddBefore ...
func (instance *YamlStepBuilder) AddBefore(before yaml.Action) *YamlStepBuilder {
	instance.Before = append(instance.Before, before)
	return instance
}

//AddAfter ...
func (instance *YamlStepBuilder) AddAfter(after yaml.Action) *YamlStepBuilder {
	instance.After = append(instance.After, after)
	return instance
}

//ToExecuteAction ...
func (instance *YamlStepBuilder) ToExecuteAction(data yaml.Action) *YamlStepBuilder {
	instance.Action = data
	return instance
}

//WithAssertion ...
func (instance *YamlStepBuilder) WithAssertion(data map[string]interface{}) *YamlStepBuilder {
	instance.Assertions = append(instance.Assertions, data)
	return instance
}

//WithExtractor ...
func (instance *YamlStepBuilder) WithExtractor(data map[string]interface{}) *YamlStepBuilder {
	instance.Extractors = append(instance.Extractors, data)
	return instance
}
