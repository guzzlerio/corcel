package processor

import "ci.guzzler.io/guzzler/corcel/core"

//JobStream ...
type JobStream interface {
	HasNext() bool
	Next() core.Job
	Reset()
	Progress() int
	Size() int
}
