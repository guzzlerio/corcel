package processor

//StepRevolvingStream ...
type StepRevolvingStream struct {
	stream StepStream
}

//HasNext ...
func (instance StepRevolvingStream) HasNext() bool {
	return true
}

//Next ...
func (instance StepRevolvingStream) Next() Step {
	if !instance.stream.HasNext() {
		instance.stream.Reset()
	}
	return instance.stream.Next()
}

//Reset ...
func (instance StepRevolvingStream) Reset() {
	instance.stream.Reset()
}

//Progress ...
func (instance StepRevolvingStream) Progress() int {
	return instance.stream.Progress()
}

//Size ...
func (instance StepRevolvingStream) Size() int {
	return instance.stream.Size()
}
