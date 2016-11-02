package report

import (
	"encoding/json"

	"github.com/hoisie/mustache"
)

//RenderHistogram ...
func RenderHistogram(node UrnComposite, times []int64) string {
	histogramLayout, _ := Asset("data/histogram.mustache")

	values := node.Value.(map[string][]float64)

	keys := []string{"x"}
	for key := range values {
		keys = append(keys, key)
	}

	data := [][]float64{}

	for i := 0; i < len(times); i++ {
		lineData := []float64{float64(times[i] / 1000 / 1000 / 1000)}
		for _, value := range values {
			lineData = append(lineData, value[i])
		}
		data = append(data, lineData)
	}

	jsonValues, _ := json.Marshal(data)
	jsonLabels, _ := json.Marshal(keys)

	return mustache.Render(string(histogramLayout), map[string]interface{}{
		"data":   string(jsonValues),
		"labels": string(jsonLabels),
	})

}
