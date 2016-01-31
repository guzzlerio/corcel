package processor

import (
	"math"
	"time"
)

//StepDurationStream ...
type StepDurationStream struct {
	stream   StepStream
	start    time.Time
	duration time.Duration
}

//CreateStepDurationStream ...
func CreateStepDurationStream(stream StepStream, duration time.Duration) *StepDurationStream {
	return &StepDurationStream{
		stream:   stream,
		duration: duration,
	}
}

//HasNext ...
func (instance *StepDurationStream) HasNext() bool {
	if instance.start.IsZero() {
		return true
	}
	return time.Since(instance.start) < instance.duration
}

//Next ...
func (instance *StepDurationStream) Next() Step {
	if instance.start.IsZero() {
		instance.start = time.Now()
	}
	if !instance.stream.HasNext() {
		instance.stream.Reset()
	}
	return instance.stream.Next()
}

//Reset ...
func (instance *StepDurationStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance *StepDurationStream) Progress() int {
	current := (float64(time.Since(instance.start).Nanoseconds()) / float64(instance.Size()))
	return int(math.Ceil(current * 100))
}

//Size ...
func (instance *StepDurationStream) Size() int {
	return int(instance.duration.Nanoseconds())
}
