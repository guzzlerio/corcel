package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"gopkg.in/alecthomas/kingpin.v2"
	yamlv2 "gopkg.in/yaml.v2"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/statistics"
)

// ServerCommand ...
type ServerCommand struct {
	Port     int
	registry *core.Registry
}

func NewServerCommand(app *kingpin.Application, registry *core.Registry) {
	c := &ServerCommand{
		registry: registry,
	}
	server := app.Command("server", "Start HTTP server").Action(c.run)
	server.Flag("port", "Port").Default("54332").IntVar(&c.Port)
}

func (instance *ServerCommand) run(c *kingpin.ParseContext) error {
	// have access to c.registry
	//Start HTTP Server
	// construct HTTP Host
	// Start HTTP Host from cmd options
	fmt.Printf("Would now be starting the HTTP server on %v\n", instance.Port)
	return nil
}

// RunCommand ...
type RunCommand struct {
	Config   *config.Configuration
	registry *core.Registry
}

func NewRunCommand(app *kingpin.Application, registry *core.Registry) {
	configuration := &config.Configuration{}
	c := &RunCommand{
		Config:   configuration,
		registry: registry,
	}
	run := app.Command("run", "Execute performance test thing").Action(c.run)
	run.Arg("file", "Corcel file contains URLs or an ExecutionPlan (see the --plan argument)").Required().StringVar(&configuration.FilePath)
	run.Flag("summary", "Output summary to STDOUT").BoolVar(&configuration.Summary)
	run.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").Default("0s").DurationVar(&configuration.Duration)
	run.Flag("wait-time", "Time to wait between each execution").Default("0s").DurationVar(&configuration.WaitTime)
	run.Flag("workers", "The number of workers to execute the requests").IntVar(&configuration.Workers)
	run.Flag("random", "Select the url at random for each execution").BoolVar(&configuration.Random)
	run.Flag("plan", "Indicate that the corcel file is an ExecutionPlan").BoolVar(&configuration.Plan)
	run.Flag("verbose", "verbosity").Short('v').Action(config.Counter).Bool()
	run.Flag("progress", "Progress reporter").EnumVar(&configuration.Progress, "bar", "logo", "none")
}

func (instance *RunCommand) run(c *kingpin.ParseContext) error {
	configuration, err := config.ParseConfiguration(instance.Config)
	if err != nil {
		errormanager.Log(err)
		os.Exit(1)
	}
	logger.ConfigureLogging(configuration)

	host := NewConsoleHost(configuration, *instance.registry)
	id, _ := host.Control.Start(configuration) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	generateExecutionOutput("./output.yml", output)

	addExecutionToHistory("./history.yml", output)

	if configuration.Summary {
		outputSummary(output)
	}
	return nil
}

func generateExecutionOutput(file string, output statistics.AggregatorSnapShot) {
	outputPath, err := filepath.Abs(file)
	check(err)
	yamlOutput, err := yamlv2.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

func addExecutionToHistory(file string, output statistics.AggregatorSnapShot) {

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

func outputSummary(snapshot statistics.AggregatorSnapShot) {
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

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}