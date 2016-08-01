package core

import "fmt"

//Job ...
type Job struct {
	ID      int
	Name    string
	Steps   []Step
	Context map[string]interface{}
}

//CreateStep ...
func (instance Job) CreateStep() Step {
	return Step{
		ID:         len(instance.Steps),
		Name:       fmt.Sprintf("Step #%v", len(instance.Steps)+1),
		JobID:      instance.ID,
		Assertions: []Assertion{},
		Extractors: []Extractor{},
	}
}

//AddStep ...
func (instance Job) AddStep(step Step) Job {
	steps := append(instance.Steps, step)
	instance.Steps = steps
	return instance
}
