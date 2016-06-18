package processor

import (
	"math"

	"ci.guzzler.io/guzzler/corcel/core"
)

//JobSequentialStream ...
type JobSequentialStream struct {
	items    []core.Job
	position int
}

//CreateJobSequentialStream ...
func CreateJobSequentialStream(items []core.Job) *JobSequentialStream {
	return &JobSequentialStream{
		items:    items,
		position: 0,
	}
}

//HasNext ...
func (instance *JobSequentialStream) HasNext() bool {
	return instance.position < instance.Size()
}

//Next ...
func (instance *JobSequentialStream) Next() core.Job {
	element := instance.items[instance.position]
	instance.position++
	return element
}

//Reset ...
func (instance *JobSequentialStream) Reset() {
	instance.position = 0
}

//Progress ...
func (instance *JobSequentialStream) Progress() int {
	current := float64(instance.position) / float64(instance.Size())
	return int(math.Floor(current * 100))
}

//Size ...
func (instance *JobSequentialStream) Size() int {
	return len(instance.items)
}
