package core

import "fmt"

//Job ...
type Job struct {
	ID         int
	Name       string
	Steps      []Step
	Context    map[string]interface{}
	nextStepID int
}

//CreateStep ...
func (instance Job) CreateStep() Step {
	return Step{
		Name:       fmt.Sprintf("Step #%v", instance.nextStepID+1),
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
	instance.nextStepID = instance.nextStepID + 1
	return instance
}
