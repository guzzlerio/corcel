package test

import "ci.guzzler.io/guzzler/corcel/serialisation/yaml"

//YamlJobBuilder ...
type YamlJobBuilder struct {
	StepBuilders []*YamlStepBuilder
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
	job := yaml.ExecutionJob{}
	for _, stepBuilder := range instance.StepBuilders {
		step := stepBuilder.Build()
		job.Steps = append(job.Steps, step)
	}

	return job
}

//CreateStep ...
func (instance *YamlJobBuilder) CreateStep() *YamlJobBuilder {
	builder := CreateStepBuilder()
	instance.StepBuilders = append(instance.StepBuilders, builder)
	return instance
}

//ToExecuteAction ...
func (instance *YamlJobBuilder) ToExecuteAction(data map[string]interface{}) *YamlJobBuilder {
	stepBuilder := instance.CurrentStepBuilder()
	stepBuilder.Action = data
	return instance
}

//WithAssertion ...
func (instance *YamlJobBuilder) WithAssertion(data map[string]interface{}) *YamlJobBuilder {
	stepBuilder := instance.CurrentStepBuilder()
	stepBuilder.Assertions = append(stepBuilder.Assertions, data)
	return instance
}

//WithExtractor ...
func (instance *YamlJobBuilder) WithExtractor(data map[string]interface{}) *YamlJobBuilder {
	stepBuilder := instance.CurrentStepBuilder()
	stepBuilder.Extractors = append(stepBuilder.Extractors, data)
	return instance
}
