package extractors

import (
	"regexp"

	"github.com/guzzlerio/corcel/core"
)

//RegexExtractor ...
type RegexExtractor struct {
	Name  string
	Key   string
	Match string
	Scope string
}

//Extract ...
//GOOD Resource
//https://github.com/StefanSchroeder/Golang-Regex-Tutorial/blob/master/01-chapter2.markdown
func (instance RegexExtractor) Extract(result core.ExecutionResult) core.ExtractionResult {
	target := result[instance.Key].(string)
	re, err := regexp.Compile(instance.Match)
	if err != nil {
		return core.ExtractionResult{}
	}
	res := re.FindString(target)

	return core.ExtractionResult{
		instance.Name: res,
		"scope":       instance.Scope,
	}
}
