package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hoisie/mustache"

	"github.com/guzzlerio/corcel/statistics"
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

func float64ArrayToStringArray(values []float64) []string {
	stringValues := []string{}
	for _, value := range values {
		stringValues = append(stringValues, fmt.Sprintf("%.0f", value))
	}

	return stringValues
}

func float64RunningSum(values []float64) []float64 {
	result := []float64{}
	for _, value := range values {
		if len(result) == 0 {
			result = append(result, value)
		} else {
			result = append(result, value+result[len(result)-1])
		}
	}

	return result
}

func int64ArrayToStringArray(values []int64) []string {
	stringValues := []string{}
	for _, value := range values {
		stringValues = append(stringValues, fmt.Sprintf("%d", value))
	}

	return stringValues
}

func int64RunningSum(values []int64) []int64 {
	result := []int64{}
	for _, value := range values {
		if len(result) == 0 {
			result = append(result, value)
		} else {
			result = append(result, value+result[len(result)-1])
		}
	}

	return result
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

	executionSummary := statistics.CreateSummary(output)

	masterLayout, _ := Asset("data/corcel.layout.mustache.html")
	summaryLayout, _ := Asset("data/summary.mustache")
	summaryGraphs, _ := Asset("data/summary_graphs.mustache")
	jsRenderLayout, _ := Asset("data/render.mustache")

	layout := mustache.Render(string(summaryLayout), map[string]interface{}{
		"throughput":         fmt.Sprintf("%.0f req/s", executionSummary.Throughput),
		"total_requests":     fmt.Sprintf("%.0f", executionSummary.TotalRequests),
		"number_of_errors":   fmt.Sprintf("%.0f", executionSummary.TotalErrors),
		"availability":       fmt.Sprintf("%.4f %%", executionSummary.Availability),
		"bytes_sent":         humanize.Bytes(uint64(executionSummary.Bytes.Sent.Total)),
		"bytes_received":     humanize.Bytes(uint64(executionSummary.Bytes.Received.Total)),
		"min_latency":        fmt.Sprintf("%.0f ms", executionSummary.MinResponseTime),
		"mean_latency":       fmt.Sprintf("%.0f ms", executionSummary.MeanResponseTime),
		"max_latency":        fmt.Sprintf("%.0f ms", executionSummary.MaxResponseTime),
		"max_bytes_received": humanize.Bytes(uint64(executionSummary.Bytes.Received.Max)),
		"min_bytes_received": humanize.Bytes(uint64(executionSummary.Bytes.Received.Min)),
		"max_bytes_sent":     humanize.Bytes(uint64(executionSummary.Bytes.Sent.Max)),
		"min_bytes_sent":     humanize.Bytes(uint64(executionSummary.Bytes.Sent.Min)),
	})

	throughputValues := float64ArrayToStringArray(output.Meters["urn:action:meter:throughput"]["rateMean"])
	sentValues := int64ArrayToStringArray(int64RunningSum(output.Counters["urn:action:counter:bytes:sent"]))
	receivedValues := int64ArrayToStringArray(int64RunningSum(output.Counters["urn:action:counter:bytes:received"]))
	requestsValues := float64ArrayToStringArray(output.Meters["urn:action:meter:throughput"]["count"])
	errorValues := int64ArrayToStringArray(output.Counters["urn:action:counter:error"])

	minLatency := int64ArrayToStringArray(output.Histograms["urn:action:histogram:duration"]["min"])
	maxLatency := int64ArrayToStringArray(output.Histograms["urn:action:histogram:duration"]["max"])
	meanLatency := int64ArrayToStringArray(output.Histograms["urn:action:histogram:duration"]["mean"])
	stdDevLatency := int64ArrayToStringArray(output.Histograms["urn:action:histogram:duration"]["stddev"])

	maxBytesSentValues := int64ArrayToStringArray(output.Histograms["urn:action:histogram:bytes:sent"]["max"])
	minBytesSentValues := int64ArrayToStringArray(output.Histograms["urn:action:histogram:bytes:sent"]["min"])

	maxBytesReceivedValues := int64ArrayToStringArray(output.Histograms["urn:action:histogram:bytes:received"]["max"])
	minBytesReceivedValues := int64ArrayToStringArray(output.Histograms["urn:action:histogram:bytes:received"]["min"])

	layout += mustache.Render(string(summaryGraphs), map[string]interface{}{
		"throughput":         strings.Join(throughputValues, ","),
		"bytes_sent":         strings.Join(sentValues, ","),
		"bytes_received":     strings.Join(receivedValues, ","),
		"total_requests":     strings.Join(requestsValues, ","),
		"errors":             strings.Join(errorValues, ","),
		"min_latency":        strings.Join(minLatency, ","),
		"max_latency":        strings.Join(maxLatency, ","),
		"mean_latency":       strings.Join(meanLatency, ","),
		"stddev_latency":     strings.Join(stdDevLatency, ","),
		"max_bytes_sent":     strings.Join(maxBytesSentValues, ","),
		"min_bytes_sent":     strings.Join(minBytesSentValues, ","),
		"max_bytes_received": strings.Join(maxBytesReceivedValues, ","),
		"min_bytes_received": strings.Join(minBytesReceivedValues, ","),
	})

	layout += mustache.Render(string(jsRenderLayout), nil)

	layout = mustache.RenderInLayout(layout, string(masterLayout), nil)

	ioutil.WriteFile("corcel-report.html", []byte(layout), 0644)
}
