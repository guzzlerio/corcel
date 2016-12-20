package yaml

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/utils"

	"github.com/satori/go.uuid"
	yamlFormat "gopkg.in/yaml.v2"
)

//PlanBuilder ...
type PlanBuilder struct {
	Iterations      int
	Random          bool
	NumberOfWorkers int
	WaitTime        string
	Duration        string
	JobBuilders     []*JobBuilder
	Context         map[string]interface{}
	Before          []Action
	After           []Action
}

//NewPlanBuilder ...
func NewPlanBuilder() *PlanBuilder {
	return &PlanBuilder{
		Iterations:      0,
		Random:          false,
		NumberOfWorkers: 1,
		Duration:        "0s",
		WaitTime:        "0s",
		JobBuilders:     []*JobBuilder{},
		Context:         map[string]interface{}{},
	}
}

//SetIterations ...
func (instance *PlanBuilder) SetIterations(value int) *PlanBuilder {
	instance.Iterations = value
	return instance
}

//SetRandom ...
func (instance *PlanBuilder) SetRandom(value bool) *PlanBuilder {
	instance.Random = value
	return instance
}

//SetDuration ...
func (instance *PlanBuilder) SetDuration(value string) *PlanBuilder {
	instance.Duration = value
	return instance
}

//SetWorkers ...
func (instance *PlanBuilder) SetWorkers(value int) *PlanBuilder {
	if value <= 0 {
		panic("Numbers of workers must be greater than 0")
	}
	instance.NumberOfWorkers = value
	return instance
}

//SetWaitTime ...
func (instance *PlanBuilder) SetWaitTime(value string) *PlanBuilder {
	_, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	instance.WaitTime = value
	return instance
}

//WithContext ...
func (instance *PlanBuilder) WithContext(context map[string]interface{}) *PlanBuilder {
	instance.Context = context
	return instance
}

//AddBefore ...
func (instance *PlanBuilder) AddBefore(before Action) *PlanBuilder {
	instance.Before = append(instance.Before, before)
	return instance
}

//AddAfter ...
func (instance *PlanBuilder) AddAfter(after Action) *PlanBuilder {
	instance.After = append(instance.After, after)
	return instance
}

//CreateJob ...
func (instance *PlanBuilder) CreateJob(arg ...string) *JobBuilder {
	var name string
	if len(arg) == 0 {
		name = ""
	} else {
		name = arg[0]
	}
	builder := NewJobBuilder(name)
	instance.JobBuilders = append(instance.JobBuilders, builder)
	return builder
}

//Build ...
func (instance *PlanBuilder) Build() (*os.File, error) {

	outputBasePath := "/tmp/corcel/plans"
	//FIXME ignored error output from MkdirAll
	os.MkdirAll(outputBasePath, 0777)

	plan := ExecutionPlan{
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
func (instance PlanBuilder) HTTPAction() http.RequestBuilder {
	return http.NewHTTPRequestBuilder()
}
