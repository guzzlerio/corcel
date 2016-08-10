package core

//Job ...
type Job struct {
	ID         int
	Name       string
	Steps      []Step
	Context    map[string]interface{}
	Before     []Action
	After      []Action
	nextStepID int
}

//CreateStep ...
func (instance Job) CreateStep() Step {
	return Step{
		ID:         instance.nextStepID,
		JobID:      instance.ID,
		Assertions: []Assertion{},
		Extractors: []Extractor{},
	}
}

//AddStep ...
func (instance Job) AddStep(step Step) Job {
	steps := append(instance.Steps, step)
	instance.Steps = steps
	instance.nextStepID++
	return instance
}
