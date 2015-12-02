package main

import (
	"github.com/gosuri/uiprogress"
)

type ProgressBar interface {
	Set(progress int) error
}

type ConsoleProgressBar struct {
	bar ProgressBar
	config *Configuration
}

func NewProgressBar(size int, config *Configuration) *ConsoleProgressBar{
	var bar ProgressBar
	switch config.Progress{
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

func (b* ConsoleProgressBar) Set(progress int) error {
	return b.bar.Set(progress)
}

type NullProgress struct { }

func (b *NullProgress) Set(progress int) error {
	return nil
}
