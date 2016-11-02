package processor

//ProgressBar ...
type ProgressBar interface {
	Set(progress int) error
}

