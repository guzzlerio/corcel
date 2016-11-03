package processor

import (
	"math/rand"
	"time"

	"github.com/guzzlerio/corcel/core"
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
