package test

import (
	"io/ioutil"
	"os"
	"time"

	"ci.guzzler.io/guzzler/corcel/serialisation/yaml"
	"ci.guzzler.io/guzzler/corcel/utils"

	yamlFormat "gopkg.in/yaml.v2"
)

//YamlPlanBuilder ...
type YamlPlanBuilder struct {
	Random          bool
	NumberOfWorkers int
	WaitTime        string
	Duration        string
	JobBuilders     []*YamlJobBuilder
	Context         map[string]interface{}
}

//NewYamlPlanBuilder ...
func NewYamlPlanBuilder() *YamlPlanBuilder {
	return &YamlPlanBuilder{
		Random:          false,
		NumberOfWorkers: 1,
		Duration:        "0s",
		WaitTime:        "0s",
		JobBuilders:     []*YamlJobBuilder{},
		Context:         map[string]interface{}{},
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

//WithContext ...
func (instance *YamlPlanBuilder) WithContext(context map[string]interface{}) *YamlPlanBuilder {
	instance.Context = context
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
	plan := yaml.ExecutionPlan{
		Random:   instance.Random,
		Workers:  instance.NumberOfWorkers,
		WaitTime: instance.WaitTime,
		Duration: instance.Duration,
		Context:  instance.Context,
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
