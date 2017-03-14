package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
)

//RegexExtractorParser ...
type RegexExtractorParser struct{}

//Parse ...
func (instance RegexExtractorParser) Parse(input map[string]interface{}) (core.Extractor, error) {

	if _, ok := input["name"]; !ok {
		return extractors.RegexExtractor{}, errors.New("name not set")
	}

	if _, ok := input["key"]; !ok {
		return extractors.RegexExtractor{}, errors.New("key not set")
	}

	if _, ok := input["match"]; !ok {
		return extractors.RegexExtractor{}, errors.New("match not set")
	}

	extractor := extractors.RegexExtractor{
		Name:  input["name"].(string),
		Key:   input["key"].(string),
		Match: input["match"].(string),
		Scope: core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor, nil
}

//Key ...
func (instance RegexExtractorParser) Key() string {
	return "RegexExtractor"
}
