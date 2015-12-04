package main

import (
	"github.com/gosuri/uiprogress"
)

//ProgressBar ...
type ProgressBar interface {
	Set(progress int) error
}

//ConsoleProgressBar ...
type ConsoleProgressBar struct {
	bar    ProgressBar
	config *Configuration
}

//NewProgressBar ...
func NewProgressBar(size int, config *Configuration) *ConsoleProgressBar {
	var bar ProgressBar
	switch config.Progress {
	case "bar":
		uiprogress.Start()
		bar = uiprogress.AddBar(size).AppendCompleted()
	case "none":
		bar = &NullProgress{}
	default:
		bar = NewLogoProgress()
	}
	return &ConsoleProgressBar{bar, config}
}

//Set ...
func (b *ConsoleProgressBar) Set(progress int) error {
	return b.bar.Set(progress)
}

//NullProgress ...
type NullProgress struct{}

//Set ...
func (b *NullProgress) Set(progress int) error {
	return nil
}
