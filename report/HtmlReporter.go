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
	registry.Add("histogram", RenderHistogram)
	layout := ""

	masterLayout, _ := Asset("data/corcel.layout.mustache.html")
	renderedComposite := composite.Render(registry, output.Times)
	connectors := composite.Connectors()

	subModel := []map[string]string{}
	for index, connector := range connectors {
		if index == 0 {
			subModel = append(subModel, map[string]string{
				"class": "active",
				"name":  connector,
			})
		} else {
			subModel = append(subModel, map[string]string{
				"name": connector,
			})
		}
	}
	model := map[string]interface{}{"tabs": subModel}

	layout = mustache.RenderInLayout(renderedComposite, string(masterLayout), model)

	ioutil.WriteFile("corcel-report.html", []byte(layout), 0644)
}
