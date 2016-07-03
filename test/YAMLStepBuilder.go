package test

import "ci.guzzler.io/guzzler/corcel/serialisation/yaml"

//YamlStepBuilder ...
type YamlStepBuilder struct {
	Action     map[string]interface{}
	Assertions []map[string]interface{}
	Extractors []map[string]interface{}
}

//Build ...
func (instance *YamlStepBuilder) Build() yaml.ExecutionStep {
	step := yaml.ExecutionStep{}
	step.Action = instance.Action
	step.Assertions = instance.Assertions
	step.Extractors = instance.Extractors
	return step
}

//CreateStepBuilder ...
func CreateStepBuilder() *YamlStepBuilder {
	builder := &YamlStepBuilder{
		Action:     map[string]interface{}{},
		Assertions: []map[string]interface{}{},
	}
	return builder
}
