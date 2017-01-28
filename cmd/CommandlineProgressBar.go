package cmd

import (
	"github.com/gosuri/uiprogress"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/processor"
)

//ConsoleProgressBar ...
type ConsoleProgressBar struct {
	bar    processor.ProgressBar
	config *config.Configuration
}

//NewProgressBar ...
func NewProgressBar(size int, config *config.Configuration) *ConsoleProgressBar {
	var bar processor.ProgressBar
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
func (this *ConsoleProgressBar) Set(progress int) error {
	return this.bar.Set(progress)
}

//NullProgress ...
type NullProgress struct{}

//Set ...
func (this *NullProgress) Set(progress int) error {
	return nil
}
