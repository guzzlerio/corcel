package test

import (
	"io/ioutil"
	"os"
	"time"

	"ci.guzzler.io/guzzler/corcel/processor"
	"ci.guzzler.io/guzzler/corcel/utils"

	"gopkg.in/yaml.v2"
)

//YamlPlanBuilder ...
type YamlPlanBuilder struct {
	Random          bool
	NumberOfWorkers int
	WaitTime        string
	Duration        string
	JobBuilders     []*YamlJobBuilder
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
	plan := processor.YamlExecutionPlan{
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
	contents, err := yaml.Marshal(&plan)
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
func (instance *YamlJobBuilder) Build() processor.YamlExecutionJob {
	job := processor.YamlExecutionJob{}
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
func (instance *YamlStepBuilder) Build() processor.YamlExecutionStep {
	step := processor.YamlExecutionStep{}
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
