package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	yamlv2 "github.com/ghodss/yaml"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/report"
	"github.com/guzzlerio/corcel/statistics"
)

// ServerCommand ...
type ServerCommand struct {
	Port     int
	registry *core.Registry
}

//NewServerCommand ...
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

//NewRunCommand ...
func NewRunCommand(app *kingpin.Application, registry *core.Registry) {
	configuration := &config.Configuration{}
	c := &RunCommand{
		Config:   configuration,
		registry: registry,
	}
	run := app.Command("run", "Execute performance test thing").Action(c.run)
	run.Arg("file", "Corcel file contains URLs or an ExecutionPlan (see the --plan argument)").Required().StringVar(&configuration.FilePath)
	run.Flag("summary", "Output summary to STDOUT").BoolVar(&configuration.Summary)
	run.Flag("iterations", "The number of iterations to run").Short('i').Default("0").IntVar(&configuration.Iterations)
	run.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").Short('d').Default("0s").DurationVar(&configuration.Duration)
	run.Flag("wait-time", "Time to wait between each execution").Default("0s").Short('t').DurationVar(&configuration.WaitTime)
	run.Flag("workers", "The number of workers to execute the requests").Short('w').IntVar(&configuration.Workers)
	run.Flag("random", "Select the url at random for each execution").Short('r').BoolVar(&configuration.Random)
	run.Flag("plan", "Indicate that the corcel file is an ExecutionPlan").BoolVar(&configuration.Plan)
	run.Flag("verbose", "verbosity").Short('v').Action(config.Counter).Bool()
	run.Flag("progress", "Progress reporter").EnumVar(&configuration.Progress, "bar", "logo", "none")
}

func (instance *RunCommand) run(c *kingpin.ParseContext) error {
	configuration, err := config.ParseConfiguration(instance.Config)
	if err != nil {
		errormanager.Log(err)
		panic("REPLACE ME THIS IS TEMP")
		os.Exit(1)
	}
	logger.ConfigureLogging(configuration)

	host := NewConsoleHost(configuration, *instance.registry)
	id, _ := host.Control.Start(configuration) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	generateExecutionOutput("./output.yml", output)

	addExecutionToHistory("./history.yml", output)

	reporter := report.CreateHTMLReporter()
	reporter.Generate(output)

	if configuration.Summary {
		summary := statistics.CreateSummary(snapshot)
		configuration.SummaryBuilder.Write(summary)
		// outputSummary(output)
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

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}
