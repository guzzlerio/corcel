package extractors

import "github.com/guzzlerio/corcel/core"

//KeyValueExtractor ...
type KeyValueExtractor struct {
	Name  string
	Key   string
	Scope string
}

//Extract ...
func (instance KeyValueExtractor) Extract(result core.ExecutionResult) core.ExtractionResult {

	var extractionResult = core.ExtractionResult{
		"scope": instance.Scope,
	}
	if target, ok := result[instance.Key]; ok {
		extractionResult[instance.Name] = target
	}

	return extractionResult

}
