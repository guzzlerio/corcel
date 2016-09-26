package http

import (
	"net/http"

	"ci.guzzler.io/guzzler/corcel/core"
)

//YamlHTTPRequestParser ...
type YamlHTTPRequestParser struct{}

//Parse ...
func (instance YamlHTTPRequestParser) Parse(input map[string]interface{}) core.Action {
	action := CreateAction()
	action.URL = input["url"].(string)
	action.Method = input["method"].(string)
	action.Headers = http.Header{}
	for key, value := range input["httpHeaders"].(map[interface{}]interface{}) {
		action.Headers.Set(key.(string), value.(string))
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
