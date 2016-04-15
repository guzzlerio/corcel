package processor

//JobRevolvingStream ...
type JobRevolvingStream struct {
	stream JobStream
}

//HasNext ...
func (instance JobRevolvingStream) HasNext() bool {
	return true
}

//Next ...
func (instance JobRevolvingStream) Next() Job {
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
