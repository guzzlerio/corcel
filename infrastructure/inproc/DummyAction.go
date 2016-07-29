package inproc

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"ci.guzzler.io/guzzler/corcel/core"
)

//DummyAction ...
type DummyAction struct {
	LogPath     string
	LogContexts bool
	Results     map[string]interface{}
	contexts    []core.ExecutionContext
}

//Execute ...
func (instance DummyAction) Execute(context core.ExecutionContext, cancellation chan struct{}) core.ExecutionResult {
	result := core.ExecutionResult{}

	for k, v := range context {
		switch value := v.(type) {
		case string:
			for key, resultValue := range instance.Results {
				replacement := strings.Replace(resultValue.(string), k, value, -1)
				instance.Results[key] = replacement
			}
		}
	}

	for key, value := range instance.Results {
		result[key] = value
	}

	if instance.LogContexts {
		if _, err := os.Stat(instance.LogPath); os.IsNotExist(err) {
			instance.contexts = []core.ExecutionContext{}
			jsonData, _ := json.Marshal(instance.contexts)
			ioutil.WriteFile(instance.LogPath, jsonData, 0644)
		}

		data, _ := ioutil.ReadFile(instance.LogPath)
		var contexts []core.ExecutionContext
		json.Unmarshal(data, &contexts)
		instance.contexts = contexts

		instance.contexts = append(instance.contexts, context)
		jsonData, _ := json.Marshal(instance.contexts)

		ioutil.WriteFile(instance.LogPath, jsonData, 0644)
	}

	return result
}
