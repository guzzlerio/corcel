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

//RegexExtractorBuilder ...
type RegexExtractorBuilder struct {
	data map[string]interface{}
}

//XPathExtractorBuilder ...
type XPathExtractorBuilder struct {
	data map[string]interface{}
}

//JSONPathExtractorBuilder ...
type JSONPathExtractorBuilder struct {
	data map[string]interface{}
}

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

//RegexExtractor ...
func (instance YamlPlanBuilder) RegexExtractor() RegexExtractorBuilder {
	return RegexExtractorBuilder{
		data: map[string]interface{}{
			"type": "RegexExtractor",
		},
	}
}

//Name ...
func (instance RegexExtractorBuilder) Name(value string) RegexExtractorBuilder {
	instance.data["name"] = value
	return instance
}

//Key ...
func (instance RegexExtractorBuilder) Key(value string) RegexExtractorBuilder {
	instance.data["key"] = value
	return instance
}

//Match ...
func (instance RegexExtractorBuilder) Match(value string) RegexExtractorBuilder {
	instance.data["match"] = value
	return instance
}

//Scope ...
func (instance RegexExtractorBuilder) Scope(value string) RegexExtractorBuilder {
	instance.data["scope"] = value
	return instance
}

//Build ...
func (instance RegexExtractorBuilder) Build() map[string]interface{} {
	return instance.data
}

//JSONPathExtractor ...
func (instance YamlPlanBuilder) JSONPathExtractor() JSONPathExtractorBuilder {
	return JSONPathExtractorBuilder{
		data: map[string]interface{}{
			"type": "JSONPathExtractor",
		},
	}
}

//Name ...
func (instance JSONPathExtractorBuilder) Name(value string) JSONPathExtractorBuilder {
	instance.data["name"] = value
	return instance
}

//Key ...
func (instance JSONPathExtractorBuilder) Key(value string) JSONPathExtractorBuilder {
	instance.data["key"] = value
	return instance
}

//JSONPath ...
func (instance JSONPathExtractorBuilder) JSONPath(value string) JSONPathExtractorBuilder {
	instance.data["jsonpath"] = value
	return instance
}

//Scope ...
func (instance JSONPathExtractorBuilder) Scope(value string) JSONPathExtractorBuilder {
	instance.data["scope"] = value
	return instance
}

//Build ...
func (instance JSONPathExtractorBuilder) Build() map[string]interface{} {
	return instance.data
}

//XPathExtractor ...
func (instance YamlPlanBuilder) XPathExtractor() XPathExtractorBuilder {
	return XPathExtractorBuilder{
		data: map[string]interface{}{
			"type": "XPathExtractor",
		},
	}
}

//Name ...
func (instance XPathExtractorBuilder) Name(value string) XPathExtractorBuilder {
	instance.data["name"] = value
	return instance
}

//Key ...
func (instance XPathExtractorBuilder) Key(value string) XPathExtractorBuilder {
	instance.data["key"] = value
	return instance
}

//XPath ...
func (instance XPathExtractorBuilder) XPath(value string) XPathExtractorBuilder {
	instance.data["xpath"] = value
	return instance
}

//Scope ...
func (instance XPathExtractorBuilder) Scope(value string) XPathExtractorBuilder {
	instance.data["scope"] = value
	return instance
}

//Build ...
func (instance XPathExtractorBuilder) Build() map[string]interface{} {
	return instance.data
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

//HTTPRequestBuilder ...
type HTTPRequestBuilder struct {
	data map[string]interface{}
}

//Timeout ...
func (instance HTTPRequestBuilder) Timeout(value int) HTTPRequestBuilder {
	instance.data["requestTimeout"] = value
	return instance
}

//Method ...
func (instance HTTPRequestBuilder) Method(value string) HTTPRequestBuilder {
	instance.data["method"] = value
	return instance
}

//URL ...
func (instance HTTPRequestBuilder) URL(value string) HTTPRequestBuilder {
	instance.data["url"] = value
	return instance
}

//Header ...
func (instance HTTPRequestBuilder) Header(key string, value string) HTTPRequestBuilder {
	instance.data["httpHeaders"].(map[string]string)[key] = value
	return instance
}

//Body ...
func (instance HTTPRequestBuilder) Body(value string) HTTPRequestBuilder {
	instance.data["body"] = value
	return instance
}

//Build ...
func (instance HTTPRequestBuilder) Build() map[string]interface{} {
	return instance.data
}

//HTTPRequestAction ...
func (instance YamlPlanBuilder) HTTPRequestAction() HTTPRequestBuilder {
	return HTTPRequestBuilder{
		data: map[string]interface{}{
			"type":        "HttpRequest",
			"method":      "GET",
			"url":         "",
			"httpHeaders": map[string]string{},
		},
	}
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

//GreaterThanAssertion ...
func (instance YamlPlanBuilder) GreaterThanAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "GreaterThanAssertion",
		"key":      key,
		"expected": expected,
	}
}

//GreaterThanOrEqualAssertion ...
func (instance YamlPlanBuilder) GreaterThanOrEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "GreaterThanOrEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}

//LessThanAssertion ...
func (instance YamlPlanBuilder) LessThanAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "LessThanAssertion",
		"key":      key,
		"expected": expected,
	}
}

//LessThanOrEqualAssertion ...
func (instance YamlPlanBuilder) LessThanOrEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "LessThanOrEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}

//NotEmptyAssertion ...
func (instance YamlPlanBuilder) NotEmptyAssertion(key string) map[string]interface{} {
	return map[string]interface{}{
		"type": "NotEmptyAssertion",
		"key":  key,
	}
}

//NotEqualAssertion ...
func (instance YamlPlanBuilder) NotEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "NotEqualAssertion",
		"key":      key,
		"expected": expected,
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
	plan := yaml.ExecutionPlan{
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
