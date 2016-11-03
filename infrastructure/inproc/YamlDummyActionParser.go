package inproc

import "github.com/guzzlerio/corcel/core"

//YamlDummyActionParser ...
type YamlDummyActionParser struct{}

//Parse ...
func (instance YamlDummyActionParser) Parse(input map[string]interface{}) core.Action {
	results := map[string]interface{}{}
	for key, value := range input["results"].(map[interface{}]interface{}) {
		results[key.(string)] = value
	}

	var logpath string

	if input["logpath"] != nil {
		logpath = input["logpath"].(string)
	}

	return DummyAction{
		Results:     results,
		LogContexts: (input["logcontexts"] != nil),
		LogPath:     logpath,
	}
}

//Key ...
func (instance YamlDummyActionParser) Key() string {
	return "DummyAction"
}
