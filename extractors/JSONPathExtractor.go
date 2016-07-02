package extractors

import (
	"fmt"

	"github.com/NodePrime/jsonpath"

	"ci.guzzler.io/guzzler/corcel/core"
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

	paths, err := jsonpath.ParsePaths(instance.JSONPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error Parsing %v", err))
		return extractionResult
	}

	fmt.Println(fmt.Sprintf("Parsing path %s", instance.JSONPath))
	eval, err := jsonpath.EvalPathsInBytes([]byte{}, paths)
	fmt.Println(fmt.Sprintf("Eval %v", eval))

	if err != nil {
		fmt.Println(fmt.Sprintf("Error Evaluating %v", err))
		return extractionResult
	}

	res := toResultArray(eval)

	extractionResult[instance.Name] = res
	fmt.Println(fmt.Sprintf("result = %v", res))

	return extractionResult

}

func toResultArray(e *jsonpath.Eval) []jsonpath.Result {
	var vals []jsonpath.Result
	for {
		if r, ok := e.Next(); ok {
			if r != nil {
				vals = append(vals, *r)
			}
		} else {
			break
		}
	}
	return vals
}
