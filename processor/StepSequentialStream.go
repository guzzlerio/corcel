package processor

import (
	"math"

	"ci.guzzler.io/guzzler/corcel/core"
)

//StepSequentialStream ...
type StepSequentialStream struct {
	items    []core.Step
	position int
}

//CreateStepSequentialStream ...
func CreateStepSequentialStream(items []core.Step) *StepSequentialStream {
	return &StepSequentialStream{
		items:    items,
		position: 0,
	}
}

//HasNext ...
func (instance *StepSequentialStream) HasNext() bool {
	return instance.position < instance.Size()
}

//Next ...
func (instance *StepSequentialStream) Next() core.Step {
	element := instance.items[instance.position]
	instance.position++
	return element
}

//Reset ...
func (instance *StepSequentialStream) Reset() {
	instance.position = 0
}

//Progress ...
func (instance *StepSequentialStream) Progress() int {
	current := float64(instance.position) / float64(instance.Size())
	return int(math.Floor(current * 100))
}

//Size ...
func (instance *StepSequentialStream) Size() int {
	return len(instance.items)
}
