package report

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hoisie/mustache"

	"ci.guzzler.io/guzzler/corcel/statistics"
)

//GraphData ...
type GraphData struct {
	Name  string
	Value float64
	Data  [][]float64
}

//DataAsJSON ...
func (instance GraphData) DataAsJSON() string {
	json, _ := json.Marshal(instance.Data)
	return string(json)
}

//HTMLReporter ...
type HTMLReporter struct {
}

//CreateHTMLReporter ...
func CreateHTMLReporter() HTMLReporter {
	return HTMLReporter{}
}

//Generate ...
func (instance HTMLReporter) Generate(output statistics.AggregatorSnapShot) {

	composite := createNode("urn", nil)

	for key, value := range output.Counters {
		composite.AddValue(key, value)
	}

	for key, value := range output.Gauges {
		composite.AddValue(key, value)
	}

	for key, value := range output.Histograms {
		composite.AddValue(key, value)
	}

	for key, value := range output.Meters {
		composite.AddValue(key, value)
	}

	for key, value := range output.Timers {
		composite.AddValue(key, value)
	}

	registry := NewRendererRegistry()
	registry.Add("counter", RenderCounter)
	layout := ""

	masterLayout, _ := Asset("data/corcel.layout.mustache.html")

	layout = mustache.RenderInLayout(composite.Render(registry), string(masterLayout), nil)

	ioutil.WriteFile("corcel-report.html", []byte(layout), 0644)
}

/*
func (instance HTMLReporter) Generate(output statistics.AggregatorSnapShot) {

	titlesReplace := []string{
		"throughput",
		"error",
		"bytes received",
		"bytes sent",
	}

	descriptionReplace := map[string]string{
		"rate1":    "1 min",
		"rate5":    "5 min",
		"rate15":   "15 min",
		"rateMean": "Avg",
		"count":    "Total",
	}

	processName := func(input string) string {
		newName := input
		for key, value := range descriptionReplace {
			if strings.Contains(strings.ToLower(input), strings.ToLower(key)) {
				newName = value
			}
		}

		for _, value := range titlesReplace {
			if strings.Contains(strings.ToLower(input), strings.ToLower(value)) {
				newName = fmt.Sprintf("%s %s", newName, value)
			}
		}

		return strings.ToTitle(newName)
	}

	//graphLayout, _ := Asset("data/graph.mustache")
	graphsLayout, _ := Asset("data/graphs.mustache")
	layout, _ := Asset("data/corcel.layout.mustache.html")

	//errorData := createErrorData(output)

	allGraphData := []map[string]string{}

	for name, meterValues := range output.Meters {
		for statKey, statValue := range meterValues {
			graphData := createGraphData(statKey, output.Times, statValue)

			allGraphData = append(allGraphData, map[string]string{
				"name":  processName(name + ":" + graphData.Name),
				"value": strconv.FormatFloat(graphData.Value, 'f', 0, 64),
				"data":  graphData.DataAsJSON(),
			})
		}
	}

	for statKey, statValue := range output.Counters {
		floatValues := []float64{}
		for _, val := range statValue {
			floatValues = append(floatValues, float64(val))
		}

		graphData := createGraphData(statKey, output.Times, floatValues)

		allGraphData = append(allGraphData, map[string]string{
			"name":  processName(graphData.Name + "Sample"),
			"value": strconv.FormatFloat(graphData.Value, 'f', 0, 64),
			"data":  graphData.DataAsJSON(),
		})
	}

	model := map[string]interface{}{
		"meters": allGraphData,
	}

	data := mustache.RenderInLayout(string(graphsLayout), string(layout), model)

	ioutil.WriteFile("corcel-report.html", []byte(data), 0644)
}
*/
