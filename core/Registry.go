package core

//Registry ...
type Registry struct {
	AssertionParsers []ExecutionAssertionParser
	ActionParsers    []ExecutionActionParser
	ResultProcessors []ExecutionResultProcessor
	ExtractorParsers []ExecutionExtractorParser
}

//CreateRegistry ...
func CreateRegistry() Registry {
	registry := Registry{
		AssertionParsers: []ExecutionAssertionParser{},
		ActionParsers:    []ExecutionActionParser{},
		ResultProcessors: []ExecutionResultProcessor{},
		ExtractorParsers: []ExecutionExtractorParser{},
	}
	return registry
}

//AddExtractorParser ...
func (instance Registry) AddExtractorParser(parser ExecutionExtractorParser) Registry {
	parsers := append(instance.ExtractorParsers, parser)
	return Registry{
		AssertionParsers: instance.AssertionParsers,
		ActionParsers:    instance.ActionParsers,
		ResultProcessors: instance.ResultProcessors,
		ExtractorParsers: parsers,
	}
}

//AddAssertionParser ...
func (instance Registry) AddAssertionParser(parser ExecutionAssertionParser) Registry {
	parsers := append(instance.AssertionParsers, parser)
	return Registry{
		AssertionParsers: parsers,
		ActionParsers:    instance.ActionParsers,
		ResultProcessors: instance.ResultProcessors,
		ExtractorParsers: instance.ExtractorParsers,
	}
}

//AddActionParser ...
func (instance Registry) AddActionParser(parser ExecutionActionParser) Registry {
	parsers := append(instance.ActionParsers, parser)
	return Registry{
		AssertionParsers: instance.AssertionParsers,
		ActionParsers:    parsers,
		ResultProcessors: instance.ResultProcessors,
		ExtractorParsers: instance.ExtractorParsers,
	}
}

//AddResultProcessor ...
func (instance Registry) AddResultProcessor(processor ExecutionResultProcessor) Registry {
	processors := append(instance.ResultProcessors, processor)
	return Registry{
		AssertionParsers: instance.AssertionParsers,
		ActionParsers:    instance.ActionParsers,
		ResultProcessors: processors,
		ExtractorParsers: instance.ExtractorParsers,
	}
}
