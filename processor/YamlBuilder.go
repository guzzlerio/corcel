package processor

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	. "ci.guzzler.io/guzzler/corcel/utils"
)

type YamlPlanBuilder struct {
	JobBuilders []*YamlJobBuilder
}

func NewYamlPlanBuilder() *YamlPlanBuilder {
	return &YamlPlanBuilder{
		JobBuilders: []*YamlJobBuilder{},
	}
}

func (instance *YamlPlanBuilder) CreateJob() *YamlJobBuilder {
	builder := NewYamlJobBuilder()
	instance.JobBuilders = append(instance.JobBuilders, builder)
	return builder
}

func (instance *YamlPlanBuilder) Build() (*os.File, error) {
	plan := YamlExecutionPlan{}
	for _, jobBuilder := range instance.JobBuilders {
		yamlExecutionJob := jobBuilder.Build()
		plan.Jobs = append(plan.Jobs, yamlExecutionJob)
	}
	file, err := ioutil.TempFile(os.TempDir(), "yamlExecutionPlanForCorcel")
	if err != nil {
		return nil, err
	}
	defer func() {
		CheckErr(file.Close())
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

type YamlJobBuilder struct {
	StepBuilders []*YamlStepBuilder
}

func NewYamlJobBuilder() *YamlJobBuilder {
	return &YamlJobBuilder{
		StepBuilders: []*YamlStepBuilder{},
	}
}

func (instance *YamlJobBuilder) CurrentStepBuilder() *YamlStepBuilder {
	if len(instance.StepBuilders) == 0 {
		panic("no builders")
	}

	return instance.StepBuilders[len(instance.StepBuilders)-1]
}

func (instance *YamlJobBuilder) Build() YamlExecutionJob {
	job := YamlExecutionJob{}
	for _, stepBuilder := range instance.StepBuilders {
		step := stepBuilder.Build()
		job.Steps = append(job.Steps, step)
	}

	return job
}

func (instance *YamlJobBuilder) CreateStep() *YamlJobBuilder {
	builder := CreateStepBuilder()
	instance.StepBuilders = append(instance.StepBuilders, builder)
	return instance
}

type YamlStepBuilder struct {
	Action     map[string]interface{}
	Assertions []map[string]interface{}
}

func (instance *YamlStepBuilder) Build() YamlExecutionStep {
	step := YamlExecutionStep{}
	step.Action = instance.Action
	step.Assertions = instance.Assertions
	return step
}

func CreateStepBuilder() *YamlStepBuilder {
	builder := &YamlStepBuilder{
		Action:     map[string]interface{}{},
		Assertions: []map[string]interface{}{},
	}
	return builder
}

func (instance *YamlJobBuilder) Action(data map[string]interface{}) *YamlJobBuilder {
	stepBuilder := instance.CurrentStepBuilder()
	stepBuilder.Action = data
	return instance
}

func (instance *YamlJobBuilder) Assertion(data map[string]interface{}) *YamlJobBuilder {
	stepBuilder := instance.CurrentStepBuilder()
	stepBuilder.Assertions = append(stepBuilder.Assertions, data)
	return instance
}
