package inproc

import "ci.guzzler.io/guzzler/corcel/core"

//YamlDummyActionParser ...
type YamlDummyActionParser struct{}

//Parse ...
func (instance YamlDummyActionParser) Parse(input map[string]interface{}) core.Action {
	results := map[string]interface{}{}
	for key, value := range input["results"].(map[interface{}]interface{}) {
		results[key.(string)] = value
	}
	return DummyAction{
		Results: results,
	}
}

//Key ...
func (instance YamlDummyActionParser) Key() string {
	return "DummyAction"
}
