package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
)

//JSONPathExtractorParser ...
type JSONPathExtractorParser struct{}

//Parse ...
func (instance JSONPathExtractorParser) Parse(input map[string]interface{}) (core.Extractor, error) {

	if _, ok := input["name"]; !ok {
		return extractors.JSONPathExtractor{}, errors.New("name not set")
	}

	if _, ok := input["key"]; !ok {
		return extractors.JSONPathExtractor{}, errors.New("key not set")
	}

	if _, ok := input["jsonpath"]; !ok {
		return extractors.JSONPathExtractor{}, errors.New("jsonpath not set")
	}

	extractor := extractors.JSONPathExtractor{
		Name:     input["name"].(string),
		Key:      input["key"].(string),
		JSONPath: input["jsonpath"].(string),
		Scope:    core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor, nil
}

//Key ...
func (instance JSONPathExtractorParser) Key() string {
	return "JSONPathExtractor"
}
