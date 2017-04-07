package yaml

import (
	"os"
	"time"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/utils"
)

//PlanBuilder ...
type PlanBuilder struct {
	Name            string
	Iterations      int
	Random          bool
	NumberOfWorkers int
	WaitTime        string
	Duration        string
	JobBuilders     []*JobBuilder
	Context         core.ExecutionContext
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

//WithName ...
func (instance *PlanBuilder) WithName(name string) *PlanBuilder {
	instance.Name = name
	return instance
}

//WithContext ...
func (instance *PlanBuilder) WithContext(context core.ExecutionContext) *PlanBuilder {
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
func (instance *PlanBuilder) Build() ExecutionPlan {
	plan := ExecutionPlan{
		Name:       instance.Name,
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
	return plan
}

//BuildAndSave ...
//TODO deprecate this
func (instance *PlanBuilder) BuildAndSave() (*os.File, error) {
	plan := instance.Build()
	outputBasePath := "/tmp/corcel/plans"
	return utils.MarshalYamlToFile(outputBasePath, plan)
}

//HTTPAction ...
func (instance PlanBuilder) HTTPAction() RequestBuilder {
	return NewHTTPRequestBuilder()
}
