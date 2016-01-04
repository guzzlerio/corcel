package cmd

import (
	"github.com/gosuri/uiprogress"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/processor"
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
func (b *ConsoleProgressBar) Set(progress int) error {
	return b.bar.Set(progress)
}

//NullProgress ...
type NullProgress struct{}

//Set ...
func (b *NullProgress) Set(progress int) error {
	return nil
}
