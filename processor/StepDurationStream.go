package processor

import "time"

//StepDurationStream ...
type StepDurationStream struct {
	stream        StepStream
	start         time.Time
	totalDuration time.Duration
}

//HasNext ...
func (instance StepDurationStream) HasNext() bool {
	if time.Since(instance.start) < instance.totalDuration {
		return instance.stream.HasNext()
	}

	return false
}

//Next ...
func (instance StepDurationStream) Next() Step {
	return instance.stream.Next()
}

//Reset ...
func (instance StepDurationStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance StepDurationStream) Progress() int {
	return instance.stream.Progress()
}

//Size ...
func (instance StepDurationStream) Size() int {
	return instance.stream.Size()
}
