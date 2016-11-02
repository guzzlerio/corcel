package processor

import (
	"math/rand"
	"time"

	"ci.guzzler.io/guzzler/corcel/core"
)

var (
	//RandomSource ...
	RandomSource = rand.NewSource(time.Now().UnixNano())
	//Random ...
	Random = rand.New(RandomSource)
)

//StepStream ...
type StepStream interface {
	HasNext() bool
	Next() core.Step
	Reset()
	Progress() int
	Size() int
}
