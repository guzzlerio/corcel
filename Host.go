package main

import (
	"ci.guzzler.io/guzzler/corcel/config"
)

// Host ...
type Host interface {
	SetControl(*Control)
}

// ConsoleHost ...
type ConsoleHost struct {
	Control Control
}

// SetControl ...
func (host *ConsoleHost) SetControl(control Control) {
	host.Control = control
}

// NewConsoleHost ...
func NewConsoleHost(config *config.Configuration) ConsoleHost {
	host := ConsoleHost{}
	bar := NewProgressBar(100, config)
	control := NewControl(bar)
	host.SetControl(control)
	return host
}
