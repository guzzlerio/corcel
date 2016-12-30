package yaml

import (
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
)

//KeyValueExtractorParser ...
type KeyValueExtractorParser struct{}

//Parse ...
func (instance KeyValueExtractorParser) Parse(input map[string]interface{}) core.Extractor {
	extractor := extractors.KeyValueExtractor{
		Name:  input["name"].(string),
		Key:   input["key"].(string),
		Scope: core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor
}

//Key ...
func (instance KeyValueExtractorParser) Key() string {
	return "KeyValueExtractor"
}
