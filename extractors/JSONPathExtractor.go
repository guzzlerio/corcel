package extractors

import (
	"encoding/json"

	"github.com/guzzlerio/jsonpath"

	"github.com/guzzlerio/corcel/core"
)

//JSONPathExtractor ...
type JSONPathExtractor struct {
	Name     string
	Key      string
	JSONPath string
	Scope    string
}

//Extract ...
func (instance JSONPathExtractor) Extract(result core.ExecutionResult) core.ExtractionResult {
	extractionResult := core.ExtractionResult{
		"scope": instance.Scope,
	}

	if data, ok := result[instance.Key]; ok {
		var json_data interface{}
		json.Unmarshal([]byte(data.(string)), &json_data)
		res, err := jsonpath.JsonPathLookup(json_data, instance.JSONPath)
		if err != nil {
			extractionResult[instance.Name] = "Unexpected error evaluating JSON Path"
		}

		extractionResult[instance.Name] = res
	}

	return extractionResult

}
