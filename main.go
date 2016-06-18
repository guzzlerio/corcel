package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	yamlv2 "gopkg.in/yaml.v2"

	"ci.guzzler.io/guzzler/corcel/cmd"
	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/infrastructure/http"
	"ci.guzzler.io/guzzler/corcel/infrastructure/inproc"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/serialisation/yaml"
	"ci.guzzler.io/guzzler/corcel/statistics"
)

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}

//GenerateExecutionOutput ...
func GenerateExecutionOutput(file string, output statistics.AggregatorSnapShot) {
	outputPath, err := filepath.Abs(file)
	check(err)
	yamlOutput, err := yamlv2.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

//AddExecutionToHistory ...
func AddExecutionToHistory(file string, output statistics.AggregatorSnapShot) {

	var summary statistics.AggregatorSnapShot

	outputPath, err := filepath.Abs(file)
	check(err)

	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		summary = *statistics.NewAggregatorSnapShot()
	} else {
		data, dataErr := ioutil.ReadFile(outputPath)
		if dataErr != nil {
			panic(dataErr)
		}
		yamlErr := yamlv2.Unmarshal(data, &summary)
		if yamlErr != nil {
			panic(yamlErr)
		}
	}
	summary.Update(output)

	yamlOutput, err := yamlv2.Marshal(&summary)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

func main() {
	logger.Initialise()
	configuration, err := config.ParseConfiguration(os.Args[1:])
	if err != nil {
		config.Usage()
		os.Exit(1)
	}

	logger.ConfigureLogging(configuration)

	registry := core.CreateRegistry().
		AddActionParser(inproc.YamlDummyActionParser{}).
		AddActionParser(http.YamlHTTPRequestParser{}).
		AddAssertionParser(yaml.YamlExactAssertionParser{}).
		AddAssertionParser(yaml.YamlEmptyAssertionParser{}).
		AddResultProcessor(http.NewHTTPExecutionResultProcessor()).
		AddResultProcessor(inproc.NewGeneralExecutionResultProcessor())

	_, err = filepath.Abs(configuration.FilePath)
	check(err)

	host := cmd.NewConsoleHost(configuration, registry)
	id, _ := host.Control.Start(configuration) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	GenerateExecutionOutput("./output.yml", output)

	AddExecutionToHistory("./history.yml", output)

	if configuration.Summary {
		OutputSummary(output)
	}
}

//OutputSummary ...
func OutputSummary(snapshot statistics.AggregatorSnapShot) {
	summary := statistics.CreateSummary(snapshot)

	top(os.Stdout)
	line(os.Stdout, "Running Time", summary.RunningTime)
	line(os.Stdout, "Throughput", fmt.Sprintf("%-.0f req/s", summary.Throughput))
	line(os.Stdout, "Total Requests", fmt.Sprintf("%-.0f", summary.TotalRequests))
	line(os.Stdout, "Number of Errors", fmt.Sprintf("%-.0f", summary.TotalErrors))
	line(os.Stdout, "Availability", fmt.Sprintf("%-.4f%%", summary.Availability))
	line(os.Stdout, "Bytes Sent", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.TotalBytesSent))))
	line(os.Stdout, "Bytes Received", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.TotalBytesReceived))))
	line(os.Stdout, "Mean Response Time", fmt.Sprintf("%.4f ms", summary.MeanResponseTime))
	line(os.Stdout, "Min Response Time", fmt.Sprintf("%.4f ms", summary.MinResponseTime))
	line(os.Stdout, "Max Response Time", fmt.Sprintf("%.4f ms", summary.MaxResponseTime))
	tail(os.Stdout)
}

func top(writer io.Writer) {
	fmt.Fprintln(writer, "╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(writer, "║                           Summary                                 ║")
	fmt.Fprintln(writer, "╠═══════════════════════════════════════════════════════════════════╣")
}

func tail(writer io.Writer) {
	fmt.Fprintln(writer, "╚═══════════════════════════════════════════════════════════════════╝")
}

func line(writer io.Writer, label string, value string) {
	fmt.Fprintf(writer, "║ %20s: %-43s ║\n", label, value)
}
