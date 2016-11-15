package yaml

//JobBuilder ...
type JobBuilder struct {
	Name         string
	StepBuilders []*StepBuilder
	Context      map[string]interface{}
	Before       []Action
	After        []Action
}

//NewJobBuilder ...
func NewJobBuilder(name string) *JobBuilder {
	return &JobBuilder{
		Name:         name,
		StepBuilders: []*StepBuilder{},
	}
}

//CurrentStepBuilder ...
func (instance *JobBuilder) CurrentStepBuilder() *StepBuilder {
	if len(instance.StepBuilders) == 0 {
		panic("no builders")
	}

	return instance.StepBuilders[len(instance.StepBuilders)-1]
}

//Build ...
func (instance *JobBuilder) Build() ExecutionJob {
	job := ExecutionJob{
		Name:    instance.Name,
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
func (instance *JobBuilder) CreateStep() *StepBuilder {
	builder := NewStepBuilder()
	instance.StepBuilders = append(instance.StepBuilders, builder)
	return builder
}

//AddBefore ...
func (instance *JobBuilder) AddBefore(before Action) *JobBuilder {
	instance.Before = append(instance.Before, before)
	return instance
}

//AddAfter ...
func (instance *JobBuilder) AddAfter(after Action) *JobBuilder {
	instance.After = append(instance.After, after)
	return instance
}

//WithContext ...
func (instance *JobBuilder) WithContext(context map[string]interface{}) *JobBuilder {
	instance.Context = context
	return instance
}
