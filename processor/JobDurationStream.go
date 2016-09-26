package processor

import (
	"math"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/core"
)

//JobDurationStream ...
type JobDurationStream struct {
	stream   JobStream
	start    time.Time
	duration time.Duration
	mutex    *sync.Mutex
}

//CreateJobDurationStream ...
func CreateJobDurationStream(stream JobStream, duration time.Duration) *JobDurationStream {
	return &JobDurationStream{
		stream:   stream,
		duration: duration,
		mutex:    &sync.Mutex{},
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
	instance.mutex.Lock()
	if instance.start.IsZero() {
		instance.start = time.Now()
	}
	if !instance.stream.HasNext() {
		instance.stream.Reset()
	}
	instance.mutex.Unlock()
	return instance.stream.Next()
}

//Reset ...
func (instance *JobDurationStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance *JobDurationStream) Progress() int {
	instance.mutex.Lock()
	current := (float64(time.Since(instance.start).Nanoseconds()) / float64(instance.Size()))
	instance.mutex.Unlock()
	return int(math.Ceil(current * 100))
}

//Size ...
func (instance *JobDurationStream) Size() int {
	return int(instance.duration.Nanoseconds())
}
