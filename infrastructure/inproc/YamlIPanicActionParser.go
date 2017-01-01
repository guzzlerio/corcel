package inproc

import "github.com/guzzlerio/corcel/core"

//YamlIPanicActionParser ...
type YamlIPanicActionParser struct{}

//Parse ...
func (instance YamlIPanicActionParser) Parse(input map[string]interface{}) core.Action {
	return IPanicAction{}
}

//Key ...
func (instance YamlIPanicActionParser) Key() string {
	return "IPanicAction"
}
