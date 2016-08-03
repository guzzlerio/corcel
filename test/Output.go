package test

import (
	"encoding/json"
	"io/ioutil"

	"ci.guzzler.io/guzzler/corcel/core"
)

//GetExecutionContexts ...
func GetExecutionContexts(path string) []core.ExecutionContext {
	data, _ := ioutil.ReadFile(path)
	var contexts []core.ExecutionContext
	json.Unmarshal(data, &contexts)
	return contexts
}
