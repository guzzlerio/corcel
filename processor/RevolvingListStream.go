package processor

//RevolvingListStream ...
type RevolvingListStream struct {
	data    []map[string]interface{}
	current int
}

//NewRevolvingListStream ...
func NewRevolvingListStream(data []map[string]interface{}) *RevolvingListStream {
	return &RevolvingListStream{
		data:    data,
		current: 0,
	}
}

//HasNext ...
func (instance *RevolvingListStream) HasNext() bool {
	return true
}

//Next ...
func (instance *RevolvingListStream) Next() map[string]interface{} {
	if instance.current == len(instance.data) {
		instance.current = 0
	}
	data := instance.data[instance.current]
	instance.current++
	return data
}

//Reset ...
func (instance *RevolvingListStream) Reset() {
	instance.current = 0
}

//Size ...
func (instance *RevolvingListStream) Size() int {
	return len(instance.data)
}
