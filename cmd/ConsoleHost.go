package cmd

import (
	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/processor"
)

// ConsoleHost ...
// This ConsoleHost should also be responsible for parsing any command line arguments
// which have been passed in
type ConsoleHost struct {
	Control processor.Control
}

// SetControl ...
func (host *ConsoleHost) SetControl(control processor.Control) {
	host.Control = control
}

// NewConsoleHost ...
func NewConsoleHost(config *config.Configuration, registry core.Registry) ConsoleHost {
	host := ConsoleHost{}
	bar := NewProgressBar(100, config)
	control := processor.NewControl(bar, registry)
	host.SetControl(control)
	return host
}
