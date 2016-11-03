package processor

import "github.com/guzzlerio/corcel/core"

//JobStream ...
type JobStream interface {
	HasNext() bool
	Next() core.Job
	Reset()
	Progress() int
	Size() int
}
