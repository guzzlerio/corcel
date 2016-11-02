package processor

//ListStream ...
type ListStream interface {
	HasNext() bool
	Next() map[string]interface{}
	Reset()
	Size() int
}
