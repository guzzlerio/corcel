package http

import (
	"fmt"
	"net/http"

	"github.com/guzzlerio/corcel/core"
)

//YamlHTTPRequestParser ...
type YamlHTTPRequestParser struct{}

//Parse ...
func (instance YamlHTTPRequestParser) Parse(input map[string]interface{}) core.Action {
	action := CreateAction()
	action.URL = input["url"].(string)
	action.Method = input["method"].(string)
	action.Headers = http.Header{}

	if value, ok := input["httpHeaders"]; ok && value != nil {
		for key, value := range input["httpHeaders"].(map[string]interface{}) {
			action.Headers.Set(key, fmt.Sprintf("%v", value))
		}
	}

	if _, ok := input["body"]; ok {
		action.Body = input["body"].(string)
	}

	return action
}

//Key ...
func (instance YamlHTTPRequestParser) Key() string {
	return "HttpRequest"
}
