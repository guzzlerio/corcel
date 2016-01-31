package processor

import (
	"math/rand"
	"time"
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
	Next() Step
	Reset()
	Progress() int
	Size() int
}
