package core

import (
	"fmt"
)

//SummaryBuilder ...
type SummaryBuilder interface {
	Write(summary ExecutionSummary)
}

//SummaryBuilderFactory
type SummaryBuilderFactory struct {
	builders map[string]SummaryBuilder
}

func NewSummaryBuilderFactory() *SummaryBuilderFactory {
	return &SummaryBuilderFactory{
		builders: make(map[string]SummaryBuilder),
	}
}

func (this *SummaryBuilderFactory) AddBuilder(format string, builder SummaryBuilder) *SummaryBuilderFactory {
	this.builders[format] = builder
	return this
}

func (this *SummaryBuilderFactory) Get(format string) SummaryBuilder {
	if builder := this.builders[format]; builder != nil {
		return builder
	}
	panic(fmt.Sprintf("Summary format '%s' is unsupported.", format))
}
