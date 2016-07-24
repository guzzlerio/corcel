package processor

import (
	"math"

	"ci.guzzler.io/guzzler/corcel/core"
)

//JobIterationStream ...
type JobIterationStream struct {
	jobCount   int
	iterations int
	count      int
	position   int
	stream     JobRevolvingStream
}

//CreateJobIterationStream ...
func CreateJobIterationStream(stream JobRevolvingStream, jobCount int, iterations int) *JobIterationStream {
	return &JobIterationStream{
		jobCount:   jobCount,
		iterations: iterations,
		count:      0,
		stream:     stream,
	}
}

//HasNext ...
func (instance *JobIterationStream) HasNext() bool {
	return instance.count < instance.iterations
}

//Next ...
func (instance *JobIterationStream) Next() core.Job {
	if instance.position == instance.jobCount-1 {
		instance.position = 0
		instance.count++
	} else {
		instance.position++
	}
	element := instance.stream.Next()
	return element
}

//Reset ...
func (instance *JobIterationStream) Reset() {
	instance.count = 0
	instance.position = 0
}

//Progress ...
func (instance *JobIterationStream) Progress() int {
	current := float64(instance.count) / float64(instance.iterations)
	return int(math.Floor(current * 100))
}

//Size ...
func (instance *JobIterationStream) Size() int {
	return instance.iterations
}
