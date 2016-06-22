package processor

import (
	"math"
	"time"

	"ci.guzzler.io/guzzler/corcel/core"
)

//JobDurationStream ...
type JobDurationStream struct {
	stream   JobStream
	start    time.Time
	duration time.Duration
}

//CreateJobDurationStream ...
func CreateJobDurationStream(stream JobStream, duration time.Duration) *JobDurationStream {
	return &JobDurationStream{
		stream:   stream,
		duration: duration,
	}
}

//HasNext ...
func (instance *JobDurationStream) HasNext() bool {
	if instance.start.IsZero() {
		return true
	}
	return time.Since(instance.start) < instance.duration
}

//Next ...
func (instance *JobDurationStream) Next() core.Job {
	if instance.start.IsZero() {
		instance.start = time.Now()
	}
	if !instance.stream.HasNext() {
		instance.stream.Reset()
	}
	return instance.stream.Next()
}

//Reset ...
func (instance *JobDurationStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance *JobDurationStream) Progress() int {
	current := (float64(time.Since(instance.start).Nanoseconds()) / float64(instance.Size()))
	return int(math.Ceil(current * 100))
}

//Size ...
func (instance *JobDurationStream) Size() int {
	return int(instance.duration.Nanoseconds())
}
