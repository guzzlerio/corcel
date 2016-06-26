package core

//Job ...
type Job struct {
	ID    int
	Name  string
	Steps []Step
}

//AddStep ...
func (instance Job) AddStep(step Step) Job {
	step.ID = len(instance.Steps)
	step.JobID = instance.ID
	steps := append(instance.Steps, step)
	instance.Steps = steps
	return instance
}
