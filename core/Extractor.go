package core

//Extractor ...
type Extractor interface {
	Extract(ExecutionResult) ExtractionResult
}
