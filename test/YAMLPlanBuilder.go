package test

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"ci.guzzler.io/guzzler/corcel/infrastructure/http"
	"ci.guzzler.io/guzzler/corcel/serialisation/yaml"
	"ci.guzzler.io/guzzler/corcel/utils"

	"github.com/satori/go.uuid"
	yamlFormat "gopkg.in/yaml.v2"
)

//YamlPlanBuilder ...
type YamlPlanBuilder struct {
	Iterations      int
	Random          bool
	NumberOfWorkers int
	WaitTime        string
	Duration        string
	JobBuilders     []*YamlJobBuilder
	Context         map[string]interface{}
	Before          []yaml.Action
	After           []yaml.Action
}

//NewYamlPlanBuilder ...
func NewYamlPlanBuilder() *YamlPlanBuilder {
	return &YamlPlanBuilder{
		Iterations:      0,
		Random:          false,
		NumberOfWorkers: 1,
		Duration:        "0s",
		WaitTime:        "0s",
		JobBuilders:     []*YamlJobBuilder{},
		Context:         map[string]interface{}{},
	}
}

//SetIterations ...
func (instance *YamlPlanBuilder) SetIterations(value int) *YamlPlanBuilder {
	instance.Iterations = value
	return instance
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

//AddBefore ...
func (instance *YamlPlanBuilder) AddBefore(before yaml.Action) *YamlPlanBuilder {
	instance.Before = append(instance.Before, before)
	return instance
}

//AddAfter ...
func (instance *YamlPlanBuilder) AddAfter(after yaml.Action) *YamlPlanBuilder {
	instance.After = append(instance.After, after)
	return instance
}

//CreateJob ...
func (instance *YamlPlanBuilder) CreateJob(arg ...string) *YamlJobBuilder {
	var name string
	if len(arg) == 0 {
		name = ""
	} else {
		name = arg[0]
	}
	builder := NewYamlJobBuilder(name)
	instance.JobBuilders = append(instance.JobBuilders, builder)
	return builder
}

//Build ...
func (instance *YamlPlanBuilder) Build() (*os.File, error) {

	outputBasePath := "/tmp/corcel/plans"
	//FIXME ignored error output from MkdirAll
	os.MkdirAll(outputBasePath, 0777)

	plan := yaml.ExecutionPlan{
		Iterations: instance.Iterations,
		Random:     instance.Random,
		Workers:    instance.NumberOfWorkers,
		WaitTime:   instance.WaitTime,
		Duration:   instance.Duration,
		Context:    instance.Context,
		Before:     instance.Before,
		After:      instance.After,
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
	// fmt.Println(string(contents[:]))
	if err != nil {
		return nil, err
	}
	//FIXME Write returns an error which is ignored...
	file.Write(contents)

	err = ioutil.WriteFile(path.Join(outputBasePath, uuid.NewV4().String()), contents, 0644)
	if err != nil {
		panic(err)
	}

	err = file.Sync()

	if err != nil {
		return nil, err
	}
	return file, nil
}

//HTTPAction ...
func (instance YamlPlanBuilder) HTTPAction() http.RequestBuilder {
	return http.NewHTTPRequestBuilder()
}
