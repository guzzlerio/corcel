package core

import "time"

//Plan ...
type Plan struct {
	Random   bool
	Workers  int
	Name     string
	WaitTime time.Duration
	Duration time.Duration
	Jobs     []Job
}

//AddJob ...
func (instance Plan) AddJob(job Job) Plan {
	job.ID = len(instance.Jobs)
	jobs := append(instance.Jobs, job)
	instance.Jobs = jobs
	return instance
}
