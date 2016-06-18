package test

import (
	"io/ioutil"
	"os"
	"time"

	"ci.guzzler.io/guzzler/corcel/serialisation/yaml"
	"ci.guzzler.io/guzzler/corcel/utils"

	yamlFormat "gopkg.in/yaml.v2"
)

//DummyActionBuilder ...
type DummyActionBuilder struct {
	data map[string]interface{}
}

//Set ...
func (instance DummyActionBuilder) Set(key string, value interface{}) DummyActionBuilder {
	instance.data["results"].(map[string]interface{})[key] = value

	return instance
}

//Build ...
func (instance DummyActionBuilder) Build() map[string]interface{} {
	return instance.data
}

//YamlPlanBuilder ...
type YamlPlanBuilder struct {
	Random          bool
	NumberOfWorkers int
	WaitTime        string
	Duration        string
	JobBuilders     []*YamlJobBuilder
}

//ExactAssertion ...
func (instance YamlPlanBuilder) ExactAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "ExactAssertion",
		"key":      key,
		"expected": expected,
	}
}

//EmptyAssertion ...
func (instance YamlPlanBuilder) EmptyAssertion(key string) map[string]interface{} {
	return map[string]interface{}{
		"type": "EmptyAssertion",
		"key":  key,
	}
}

//DummyAction ...
func (instance YamlPlanBuilder) DummyAction() DummyActionBuilder {
	return DummyActionBuilder{
		data: map[string]interface{}{
			"type":    "DummyAction",
			"results": map[string]interface{}{},
		},
	}
}

//NewYamlPlanBuilder ...
func NewYamlPlanBuilder() *YamlPlanBuilder {
	return &YamlPlanBuilder{
		Random:          false,
		NumberOfWorkers: 1,
		Duration:        "0s",
		WaitTime:        "0s",
		JobBuilders:     []*YamlJobBuilder{},
	}
}

//SetRandom ...
func (instance *YamlPlanBuilder) SetRandom(value bool) *YamlPlanBuilder {
	instance.Random = value
	return instance
}

//SetDuration ...
func (instance *YamlPlanBuilder) SetDuration(value string) *YamlPlanBuilder {
	instance.Duration = value
	return instance
}

//SetWorkers ...
func (instance *YamlPlanBuilder) SetWorkers(value int) *YamlPlanBuilder {
	if value <= 0 {
		panic("Numbers of workers must be greater than 0")
	}
	instance.NumberOfWorkers = value
	return instance
}

//SetWaitTime ...
func (instance *YamlPlanBuilder) SetWaitTime(value string) *YamlPlanBuilder {
	_, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	instance.WaitTime = value
	return instance
}

//CreateJob ...
func (instance *YamlPlanBuilder) CreateJob() *YamlJobBuilder {
	builder := NewYamlJobBuilder()
	instance.JobBuilders = append(instance.JobBuilders, builder)
	return builder
}

//Build ...
func (instance *YamlPlanBuilder) Build() (*os.File, error) {
	plan := yaml.YamlExecutionPlan{
		Random:   instance.Random,
		Workers:  instance.NumberOfWorkers,
		WaitTime: instance.WaitTime,
		Duration: instance.Duration,
	}
	for _, jobBuilder := range instance.JobBuilders {
		yamlExecutionJob := jobBuilder.Build()
		plan.Jobs = append(plan.Jobs, yamlExecutionJob)
	}
	file, err := ioutil.TempFile(os.TempDir(), "yamlExecutionPlanForCorcel")
	if err != nil {
		return nil, err
	}
	defer func() {
		utils.CheckErr(file.Close())
	}()
	contents, err := yamlFormat.Marshal(&plan)
	if err != nil {
		return nil, err
	}
	file.Write(contents)
	err = file.Sync()
	if err != nil {
		return nil, err
	}
	return file, nil
}

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
func (instance *YamlJobBuilder) Build() yaml.YamlExecutionJob {
	job := yaml.YamlExecutionJob{}
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

//YamlStepBuilder ...
type YamlStepBuilder struct {
	Action     map[string]interface{}
	Assertions []map[string]interface{}
}

//Build ...
func (instance *YamlStepBuilder) Build() yaml.YamlExecutionStep {
	step := yaml.YamlExecutionStep{}
	step.Action = instance.Action
	step.Assertions = instance.Assertions
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
