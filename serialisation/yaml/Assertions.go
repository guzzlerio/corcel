package yaml

import (
	"reflect"

	goyaml "github.com/ghodss/yaml"
)

func ShouldMatchYaml(actual interface{}, expected ...interface{}) string {
	var a interface{}
	var b interface{}

	goyaml.Unmarshal([]byte(actual.(string)), &a)
	goyaml.Unmarshal([]byte(expected[0].(string)), &b)

	if reflect.DeepEqual(a, b) {
		return ""
	}
	return "YAML does not match"
}
