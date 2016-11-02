package core

import (
	"fmt"
	"time"
)

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
}

//CreateJob ...
func (instance Plan) CreateJob() Job {
	return Job{
		Name:  fmt.Sprintf("Job #%v", len(instance.Jobs)+1),
		ID:    len(instance.Jobs),
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
	return instance
}
