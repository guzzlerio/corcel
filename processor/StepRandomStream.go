package processor

import "math"

//StepRandomStream ...
type StepRandomStream struct {
	items []Step
	count int
}

//CreateStepRandomStream ...
func CreateStepRandomStream(items []Step) *StepRandomStream {
	return &StepRandomStream{
		items: items,
		count: 0,
	}
}

//HasNext ...
func (instance *StepRandomStream) HasNext() bool {
	return instance.count < instance.Size()
}

//Next ...
func (instance *StepRandomStream) Next() Step {
	max := instance.Size() - 1
	if max == 0 {
		max = 1
	}
	randomIndex := Random.Intn(max)
	element := instance.items[randomIndex]

	instance.count++
	return element
}

//Reset ...
func (instance *StepRandomStream) Reset() {
	instance.count = 0
}

//Progress ...
func (instance *StepRandomStream) Progress() int {
	current := float64(instance.count) / float64(instance.Size())
	return int(math.Floor(current * 100))
}

//Size ...
func (instance *StepRandomStream) Size() int {
	return len(instance.items)
}
