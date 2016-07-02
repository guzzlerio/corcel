package core

import "time"

//Plan ...
type Plan struct {
	Random    bool
	Workers   int
	Name      string
	WaitTime  time.Duration
	Duration  time.Duration
	Jobs      []Job
	nextJobID int
}

//CreateJob ...
func (instance Plan) CreateJob() Job {
	return Job{
		ID:    instance.nextJobID,
		Steps: []Step{},
	}
}

//AddJob ...
func (instance Plan) AddJob(job Job) Plan {
	jobs := append(instance.Jobs, job)
	instance.Jobs = jobs
	instance.nextJobID = instance.nextJobID + 1
	return instance
}
