package core

import "time"

//Plan ...
type Plan struct {
	Iterations int
	Random     bool
	Workers    int
	Name       string
	WaitTime   time.Duration
	Duration   time.Duration
	Jobs       []Job
	Context    map[string]interface{}
	Before     []Action
	After      []Action
	nextJobID  int
}

//CreateJob ...
func (instance Plan) CreateJob() Job {
	return Job{
		ID:    instance.nextJobID,
		Steps: []Step{},
	}
}

//GetJob ...
func (instance Plan) GetJob(id int) Job {
	return instance.Jobs[id]
}

//AddJob ...
func (instance Plan) AddJob(job Job) Plan {
	jobs := append(instance.Jobs, job)
	instance.Jobs = jobs
	instance.nextJobID++
	return instance
}
