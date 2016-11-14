package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

//HashMapList ...
type HashMapList map[string][]map[string]interface{}

//NullPlan returns an empty initialized Plan
func NullPlan() Plan {
	return Plan{}
}

//DefaultPlan returns an initialized Plan with the default values
func DefaultPlan() Plan {
	return Plan{
		Workers: 1,
		Jobs:    []Job{},
		Context: map[string]interface{}{},
		Before:  []Action{},
		After:   []Action{},
	}

}

//PlanBuilder is a builder which ensures defaults are added to a Plan
type PlanBuilder struct {
	plan Plan
}

//NewPlanBuilder creates a new PlanBuilder
func NewPlanBuilder() PlanBuilder {
	return PlanBuilder{
		plan: DefaultPlan(),
	}
}

//Name sets the name of the plan
func (instance PlanBuilder) Name(value string) PlanBuilder {
	instance.plan.Name = value
	return PlanBuilder{
		plan: instance.plan,
	}
}

//Workers ...
func (instance PlanBuilder) Workers(value int) PlanBuilder {
	instance.plan.Workers = value
	return PlanBuilder{
		plan: instance.plan,
	}
}

//WaitTime ...
func (instance PlanBuilder) WaitTime(value time.Duration) PlanBuilder {
	instance.plan.WaitTime = value
	return PlanBuilder{
		plan: instance.plan,
	}
}

//Build returns the created plan
func (instance PlanBuilder) Build() Plan {
	return instance.plan
}

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

//Lists returns the configured lists for the plan
func (instance Plan) Lists() HashMapList {
	var lists = HashMapList{}

	if instance.Context["lists"] != nil {
		listKeys := instance.Context["lists"].(map[interface{}]interface{})
		for listKey, listValue := range listKeys {
			lists[listKey.(string)] = []map[string]interface{}{}
			listValueItems := listValue.([]interface{})
			for _, listValueItem := range listValueItems {
				srcData := listValueItem.(map[interface{}]interface{})
				stringKeyData := map[string]interface{}{}
				for srcKey, srcValue := range srcData {
					stringKeyData[srcKey.(string)] = srcValue
				}
				lists[listKey.(string)] = append(lists[listKey.(string)], stringKeyData)
			}
		}
	}
	return lists
}

//Clone returns a clone of the current plan
func (instance Plan) Clone() Plan {

	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	err := enc.Encode(instance)
	if err != nil {
		panic(err)
	}

	var clonedPlan Plan

	err = dec.Decode(&clonedPlan)
	if err != nil {
		panic(err)
	}

	return clonedPlan
}
