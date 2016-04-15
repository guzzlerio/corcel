package processor

//JobStream ...
type JobStream interface {
	HasNext() bool
	Next() Job
	Reset()
	Progress() int
	Size() int
}
