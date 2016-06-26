package yaml

import (
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/extractors"
)

//RegexExtractorParser ...
type RegexExtractorParser struct{}

//Parse ...
func (instance RegexExtractorParser) Parse(input map[string]interface{}) core.Extractor {
	extractor := extractors.RegexExtractor{
		Name:  input["name"].(string),
		Key:   input["key"].(string),
		Match: input["match"].(string),
		Scope: core.StepScope,
	}

	if input["scope"] != nil {
		extractor.Scope = input["scope"].(string)
	}

	return extractor
}

//Key ...
func (instance RegexExtractorParser) Key() string {
	return "RegexExtractor"
}
