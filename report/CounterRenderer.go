package report

import (
	"encoding/json"

	"github.com/hoisie/mustache"
)

//RenderCounter ...
func RenderCounter(node UrnComposite, times []int64) string {
	counterLayout, _ := Asset("data/counter.mustache")
	values := node.Value.([]int64)

	data := [][]int64{}

	keys := []string{"x", "value"}

	for i := 0; i < len(times); i++ {
		data = append(data, []int64{times[i] / 1000 / 1000 / 1000, values[i]})
	}

	jsonValues, _ := json.Marshal(data)
	jsonKeys, _ := json.Marshal(keys)

	return mustache.Render(string(counterLayout), map[string]interface{}{
		"data":   string(jsonValues),
		"labels": string(jsonKeys),
	})

}
