package main

import (
	"github.com/gosuri/uiprogress"
)

type ProgressBar struct {
	bar *uiprogress.Bar
	config *Configuration
}

func NewProgressBar(size int, config *Configuration) *ProgressBar{
	uiprogress.Start()
	bar := uiprogress.AddBar(size).AppendCompleted()
	return &ProgressBar{bar, config}
}

func (b* ProgressBar) Set(progress int) {
	b.bar.Set(progress)
}
