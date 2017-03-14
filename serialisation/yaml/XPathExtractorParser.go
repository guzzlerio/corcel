package yaml

import (
	"errors"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
)

//XPathExtractorParser ...
type XPathExtractorParser struct{}

//Parse ...
func (instance XPathExtractorParser) Parse(input map[string]interface{}) (core.Extractor, error) {

	if _, ok := input["name"]; !ok {
		return extractors.XPathExtractor{}, errors.New("name not found")
	}

	if _, ok := input["key"]; !ok {
		return extractors.XPathExtractor{}, errors.New("key not found")
	}

	if _, ok := input["xpath"]; !ok {
		return extractors.XPathExtractor{}, errors.New("xpath not found")
	}
	extractor := extractors.XPathExtractor{
		Name:  input["name"].(string),
		Key:   input["key"].(string),
		XPath: input["xpath"].(string),
		Scope: core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor, nil
}

//Key ...
func (instance XPathExtractorParser) Key() string {
	return "XPathExtractor"
}
