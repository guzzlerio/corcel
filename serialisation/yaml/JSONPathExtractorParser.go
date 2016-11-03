package yaml

import (
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
)

//JSONPathExtractorParser ...
type JSONPathExtractorParser struct{}

//Parse ...
func (instance JSONPathExtractorParser) Parse(input map[string]interface{}) core.Extractor {
	extractor := extractors.JSONPathExtractor{
		Name:     input["name"].(string),
		Key:      input["key"].(string),
		JSONPath: input["jsonpath"].(string),
		Scope:    core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor
}

//Key ...
func (instance JSONPathExtractorParser) Key() string {
	return "JSONPathExtractor"
}
