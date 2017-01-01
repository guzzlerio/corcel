package inproc

import (
	"context"

	"github.com/guzzlerio/corcel/core"
)

//IPanicAction ...
type IPanicAction struct {
}

//Execute ...
func (instance IPanicAction) Execute(ctx context.Context, executionContext core.ExecutionContext) core.ExecutionResult {
	panic("IPanicAction has caused this panic")
}
