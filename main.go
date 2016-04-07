package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	"ci.guzzler.io/guzzler/corcel/cmd"
	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/logger"
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
	yamlOutput, err := yaml.Marshal(&output)
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

	_, err = filepath.Abs(configuration.FilePath)
	check(err)

	host := cmd.NewConsoleHost(configuration)
	id, _ := host.Control.Start(configuration) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	GenerateExecutionOutput("./output.yml", output)

	if configuration.Summary {
		OutputSummary(output)
	}
}

func OutputSummary(snapshot statistics.AggregatorSnapShot) {
	top(os.Stdout)
	lastTime := time.Unix(snapshot.Times[len(snapshot.Times)-1], 0)
	firstTime := time.Unix(snapshot.Times[0], 0)
	duration := lastTime.Sub(firstTime)
	line(os.Stdout, "Running Time", duration.String())
	tail(os.Stdout)
}

/*
	line(writer, "Running Time", fmt.Sprintf("%g s", w.Output.Summary.RunningTime/1000))
	line(writer, "Throughput", fmt.Sprintf("%-v req/s", int64(w.Output.Summary.Requests.Rate)))
	line(writer, "Total Requests", fmt.Sprintf("%v", w.Output.Summary.Requests.Total))
	line(writer, "Number of Errors", fmt.Sprintf("%v", w.Output.Summary.Requests.Errors))
	line(writer, "Availability", fmt.Sprintf("%.4v%%", w.Output.Summary.Requests.Availability*100))
	line(writer, "Bytes Sent", fmt.Sprintf("%v", w.Output.Summary.Bytes.Sent.Sum))
	line(writer, "Bytes Received", fmt.Sprintf("%v", w.Output.Summary.Bytes.Received.Sum))
	if w.Output.Summary.ResponseTime.Mean > 0 {
		line(writer, "Mean Response Time", fmt.Sprintf("%.4v ms", w.Output.Summary.ResponseTime.Mean))
	} else {
		line(writer, "Mean Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Mean))
	}

	line(writer, "Min Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Min))
	line(writer, "Max Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Max))

*/

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
