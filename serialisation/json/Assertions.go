package json

import (
	"encoding/json"
	"reflect"
)

func ShouldMatchJson(actual interface{}, expected ...interface{}) string {
	var a interface{}
	var b interface{}

	json.Unmarshal([]byte(actual.(string)), &a)
	json.Unmarshal([]byte(expected[0].(string)), &b)

	if reflect.DeepEqual(a, b) {
		return ""
	}
	return "YAML does not match"
}
