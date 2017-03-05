package http

import nethttp "net/http"

type defaultsExtractor struct {
}

func (self defaultsExtractor) Extract(input map[string]interface{}) HttpActionState {

	var state = HttpActionState{
		Headers: nethttp.Header{},
		Method:  "GET",
	}
	if input == nil {
		return state
	}

	if _, ok := input["defaults"]; !ok {
		return state
	}
	var defaults = input["defaults"].(map[string]interface{})

	if _, ok := defaults["HttpAction"]; !ok {
		return state
	}
	var action = defaults["HttpAction"].(map[string]interface{})

	if _, ok := action["headers"]; ok {
		var headers = action["headers"].(map[string]interface{})

		for k, v := range headers {
			state.Headers.Add(k, v.(string))
		}
	}

	if _, ok := action["method"]; ok {
		state.Method = action["method"].(string)
	}

	if _, ok := action["body"]; ok {
		state.Body = action["body"].(string)
	}

	if _, ok := action["url"]; ok {
		state.URL = action["url"].(string)
	}

	return state
}

func NewDefaultsExtractor() defaultsExtractor {
	return defaultsExtractor{}
}
