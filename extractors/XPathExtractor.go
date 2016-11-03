package extractors

import (
	"bytes"
	"log"

	"gopkg.in/xmlpath.v2"

	"github.com/guzzlerio/corcel/core"
)

//XPathExtractor ...
type XPathExtractor struct {
	Name  string
	Key   string
	XPath string
	Scope string
}

//Extract ...
//GOOD Resource
//https://github.com/StefanSchroeder/Golang-XPath-Tutorial/blob/master/01-chapter2.markdown
func (instance XPathExtractor) Extract(result core.ExecutionResult) core.ExtractionResult {
	path := xmlpath.MustCompile(instance.XPath)
	root, err := xmlpath.Parse(bytes.NewBuffer([]byte(result[instance.Key].(string))))
	if err != nil {
		log.Fatal(err)
	}

	extractionResult := core.ExtractionResult{
		"scope": instance.Scope,
	}

	if value, ok := path.String(root); ok {
		extractionResult[instance.Name] = value
	}

	return extractionResult

}
