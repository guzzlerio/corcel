package processor

import (
	"math"

	"github.com/guzzlerio/corcel/core"
)

//JobRandomStream ...
type JobRandomStream struct {
	items []core.Job
	count int
}

//CreateJobRandomStream ...
func CreateJobRandomStream(items []core.Job) *JobRandomStream {
	return &JobRandomStream{
		items: items,
		count: 0,
	}
}

//HasNext ...
func (instance *JobRandomStream) HasNext() bool {
	return instance.count < instance.Size()
}

//Next ...
func (instance *JobRandomStream) Next() core.Job {
	max := instance.Size() - 1
	if max == 0 {
		max = 1
	}
	randomIndex := RandomMax(max)
	element := instance.items[randomIndex]

	instance.count++
	return element
}

//Reset ...
func (instance *JobRandomStream) Reset() {
	instance.count = 0
}

//Progress ...
func (instance *JobRandomStream) Progress() int {
	current := float64(instance.count) / float64(instance.Size())
	return int(math.Floor(current * 100))
}

//Size ...
func (instance *JobRandomStream) Size() int {
	return len(instance.items)
}
