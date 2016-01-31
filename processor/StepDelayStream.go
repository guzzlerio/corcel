package processor

import "time"

//StepDelayStream ...
type StepDelayStream struct {
	stream StepStream
	delay  time.Duration
}

//HasNext ...
func (instance StepDelayStream) HasNext() bool {
	return instance.stream.HasNext()
}

//Next ...
func (instance StepDelayStream) Next() Step {
	element := instance.stream.Next()
	time.Sleep(instance.delay)
	return element
}

//Reset ...
func (instance StepDelayStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance StepDelayStream) Progress() int {
	return instance.stream.Progress()
}

//Size ...
func (instance StepDelayStream) Size() int {
	return instance.stream.Size()
}
