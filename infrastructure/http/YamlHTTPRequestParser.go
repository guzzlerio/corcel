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
	state := HttpActionState{}

	if _, ok := input["url"]; ok {
		state.URL = input["url"].(string)
	}

	if _, ok := input["method"]; ok {
		state.Method = input["method"].(string)
	}
	state.Headers = http.Header{}

	if value, ok := input["headers"]; ok && value != nil {
		for key, value := range input["headers"].(map[string]interface{}) {
			state.Headers.Set(key, fmt.Sprintf("%v", value))
		}
	}

	if _, ok := input["body"]; ok {
		state.Body = input["body"].(string)
	}

	action.State = state

	return action
}

//Key ...
func (instance YamlHTTPRequestParser) Key() string {
	return "HttpRequest"
}
