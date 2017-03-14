package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
)

//KeyValueExtractorParser ...
type KeyValueExtractorParser struct{}

//Parse ...
func (instance KeyValueExtractorParser) Parse(input map[string]interface{}) (core.Extractor, error) {

	if _, ok := input["name"]; !ok {
		return extractors.KeyValueExtractor{}, errors.New("name not set")
	}

	if _, ok := input["key"]; !ok {
		return extractors.KeyValueExtractor{}, errors.New("key not set")
	}

	extractor := extractors.KeyValueExtractor{
		Name:  input["name"].(string),
		Key:   input["key"].(string),
		Scope: core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor, nil
}

//Key ...
func (instance KeyValueExtractorParser) Key() string {
	return "KeyValueExtractor"
}
