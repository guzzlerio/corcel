package core

import (
	"time"
)

//Step ...
type Step struct {
	Name       string
	Action     Action
	Assertions []Assertion
}

//Job ...
type Job struct {
	Name  string
	Steps []Step
}

//Plan ...
type Plan struct {
	Random   bool
	Workers  int
	Name     string
	WaitTime time.Duration
	Duration time.Duration
	Jobs     []Job
}
