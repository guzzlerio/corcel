package test

import "ci.guzzler.io/guzzler/corcel/serialisation/yaml"

//YamlJobBuilder ...
type YamlJobBuilder struct {
	StepBuilders []*YamlStepBuilder
	Context      map[string]interface{}
	Before       []yaml.Action
	After        []yaml.Action
}

//NewYamlJobBuilder ...
func NewYamlJobBuilder() *YamlJobBuilder {
	return &YamlJobBuilder{
		StepBuilders: []*YamlStepBuilder{},
	}
}

//CurrentStepBuilder ...
func (instance *YamlJobBuilder) CurrentStepBuilder() *YamlStepBuilder {
	if len(instance.StepBuilders) == 0 {
		panic("no builders")
	}

	return instance.StepBuilders[len(instance.StepBuilders)-1]
}

//Build ...
func (instance *YamlJobBuilder) Build() yaml.ExecutionJob {
	job := yaml.ExecutionJob{
		Name:    "test",
		Context: instance.Context,
		Before:  instance.Before,
		After:   instance.After,
	}
	for _, stepBuilder := range instance.StepBuilders {
		step := stepBuilder.Build()
		job.Steps = append(job.Steps, step)
	}

	return job
}

//CreateStep ...
func (instance *YamlJobBuilder) CreateStep() *YamlStepBuilder {
	builder := CreateStepBuilder()
	instance.StepBuilders = append(instance.StepBuilders, builder)
	return builder
}

//AddBefore ...
func (instance *YamlJobBuilder) AddBefore(before yaml.Action) *YamlJobBuilder {
	instance.Before = append(instance.Before, before)
	return instance
}

//AddAfter ...
func (instance *YamlJobBuilder) AddAfter(after yaml.Action) *YamlJobBuilder {
	instance.After = append(instance.After, after)
	return instance
}

//WithContext ...
func (instance *YamlJobBuilder) WithContext(context map[string]interface{}) *YamlJobBuilder {
	instance.Context = context
	return instance
}
