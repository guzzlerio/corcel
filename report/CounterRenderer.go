package report

import (
	"encoding/json"

	"github.com/hoisie/mustache"
)

//RenderCounter ...
func RenderCounter(node UrnComposite) string {
	counterLayout, _ := Asset("data/counter.mustache")
	values := [][]int64{node.Value.([]int64)}
	jsonValues, _ := json.Marshal(values)

	return mustache.Render(string(counterLayout), map[string]interface{}{
		"data": string(jsonValues),
	})

}
