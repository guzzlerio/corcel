package yaml

import (
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/extractors"
)

//XPathExtractorParser ...
type XPathExtractorParser struct{}

//Parse ...
func (instance XPathExtractorParser) Parse(input map[string]interface{}) core.Extractor {
	extractor := extractors.XPathExtractor{
		Name:  input["name"].(string),
		Key:   input["key"].(string),
		XPath: input["xpath"].(string),
		Scope: core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor
}

//Key ...
func (instance XPathExtractorParser) Key() string {
	return "XPathExtractor"
}
