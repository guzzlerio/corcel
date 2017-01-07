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

	//RandomMax ...
	RandomMax = func(max int) int {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(max)
	}
)

//StepStream ...
type StepStream interface {
	HasNext() bool
	Next() core.Step
	Reset()
	Progress() int
	Size() int
}
