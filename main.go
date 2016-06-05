package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
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

func createSummary(output statistics.AggregatorSnapShot) statistics.AggregatorSnapShot {

	var history = statistics.NewAggregatorSnapShot()
	for key, value := range output.Counters {
		//history.Counters[key] = []int64{value[len(value)-1]}
		history.UpdateCounter(key, value[len(value)-1])
	}
	for key, value := range output.Guages {
		//history.Guages[key] = []float64{value[len(value)-1]}
		history.UpdateGuage(key, value[len(value)-1])
	}
	for key, value := range output.Histograms {
		for subKey, subValue := range value {
			//history.Histograms[key][subKey] = []float64{subValue[len(subValue)-1]}
			history.UpdateHistogram(key, subKey, subValue[len(subValue)-1])
		}
	}
	for key, value := range output.Meters {
		for subKey, subValue := range value {
			//history.Meters[key][subKey] = []float64{subValue[len(subValue)-1]}
			history.UpdateMeter(key, subKey, subValue[len(subValue)-1])
		}
	}
	for key, value := range output.Timers {
		for subKey, subValue := range value {
			//history.Timers[key][subKey] = []float64{subValue[len(subValue)-1]}
			history.UpdateTimer(key, subKey, subValue[len(subValue)-1])
		}
	}
	//history.Times = []int64{output.Times[len(output.Times)-1]}
	history.UpdateTime(output.Times[len(output.Times)-1])
	return *history
}

func updateSummary(outputPath string, output statistics.AggregatorSnapShot) statistics.AggregatorSnapShot {
	var history statistics.AggregatorSnapShot
	data, err := ioutil.ReadFile(outputPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &history)
	if err != nil {
		panic(err)
	}
	for key, value := range output.Counters {
		//history.Counters[key] = append(history.Counters[key], value[len(value)-1])
		history.UpdateCounter(key, value[len(value)-1])
	}
	for key, value := range output.Guages {
		//history.Guages[key] = append(history.Guages[key], value[len(value)-1])
		history.UpdateGuage(key, value[len(value)-1])
	}
	for key, value := range output.Histograms {
		for subKey, subValue := range value {
			//history.Histograms[key][subKey] = append(history.Histograms[key][subKey], subValue[len(subValue)-1])
			history.UpdateHistogram(key, subKey, subValue[len(subValue)-1])
		}
	}
	for key, value := range output.Meters {
		for subKey, subValue := range value {
			//history.Meters[key][subKey] = append(history.Meters[key][subKey], subValue[len(subValue)-1])
			history.UpdateMeter(key, subKey, subValue[len(subValue)-1])
		}
	}
	for key, value := range output.Timers {
		for subKey, subValue := range value {
			//history.Timers[key][subKey] = append(history.Timers[key][subKey], subValue[len(subValue)-1])
			history.UpdateTimer(key, subKey, subValue[len(subValue)-1])
		}
	}
	//history.Times = append(history.Times, output.Times[len(output.Times)-1])
	history.UpdateTime(output.Times[len(output.Times)-1])
	return history
}

//AddExecutionToHistory ...
func AddExecutionToHistory(file string, output statistics.AggregatorSnapShot) {

	var history = output

	outputPath, err := filepath.Abs(file)
	check(err)

	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		history = createSummary(output)
	} else {
		history = updateSummary(outputPath, output)
	}

	yamlOutput, err := yaml.Marshal(&history)
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
