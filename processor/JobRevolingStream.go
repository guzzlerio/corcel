package processor

import "ci.guzzler.io/guzzler/corcel/core"

//JobRevolvingStream ...
type JobRevolvingStream struct {
	stream JobStream
}

//CreateJobRevolvingStream ...
func CreateJobRevolvingStream(stream JobStream) *JobRevolvingStream {
	return &JobRevolvingStream{
		stream: stream,
	}
}

//HasNext ...
func (instance JobRevolvingStream) HasNext() bool {
	return true
}

//Next ...
func (instance JobRevolvingStream) Next() core.Job {
	if !instance.stream.HasNext() {
		instance.stream.Reset()
	}
	return instance.stream.Next()
}

//Reset ...
func (instance JobRevolvingStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance JobRevolvingStream) Progress() int {
	return instance.stream.Progress()
}

//Size ...
func (instance JobRevolvingStream) Size() int {
	return instance.stream.Size()
}
