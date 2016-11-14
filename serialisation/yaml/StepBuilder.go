package yaml

//StepBuilder ...
type StepBuilder struct {
	Before     []Action
	Action     map[string]interface{}
	Assertions []map[string]interface{}
	Extractors []map[string]interface{}
	After      []Action
}

//Build ...
func (instance *StepBuilder) Build() ExecutionStep {
	return ExecutionStep{
		Before:     instance.Before,
		Action:     instance.Action,
		Assertions: instance.Assertions,
		Extractors: instance.Extractors,
		After:      instance.After,
	}
}

//NewStepBuilder ...
func NewStepBuilder() *StepBuilder {
	return &StepBuilder{
		Before:     []Action{},
		Action:     map[string]interface{}{},
		Assertions: []map[string]interface{}{},
		After:      []Action{},
	}
}

//AddBefore ...
func (instance *StepBuilder) AddBefore(before Action) *StepBuilder {
	instance.Before = append(instance.Before, before)
	return instance
}

//AddAfter ...
func (instance *StepBuilder) AddAfter(after Action) *StepBuilder {
	instance.After = append(instance.After, after)
	return instance
}

//ToExecuteAction ...
func (instance *StepBuilder) ToExecuteAction(data Action) *StepBuilder {
	instance.Action = data
	return instance
}

//WithAssertion ...
func (instance *StepBuilder) WithAssertion(data map[string]interface{}) *StepBuilder {
	instance.Assertions = append(instance.Assertions, data)
	return instance
}

//WithExtractor ...
func (instance *StepBuilder) WithExtractor(data map[string]interface{}) *StepBuilder {
	instance.Extractors = append(instance.Extractors, data)
	return instance
}
