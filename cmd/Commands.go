package cmd

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/report"
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

func (this *ServerCommand) run(c *kingpin.ParseContext) error {
	// have access to c.registry
	//Start HTTP Server
	// construct HTTP Host
	// Start HTTP Host from cmd options
	fmt.Printf("Would now be starting the HTTP server on %v\n", this.Port)
	return nil
}

// RunCommand ...
type RunCommand struct {
	Config          *config.Configuration
	registry        *core.Registry
	summaryBuilders *core.SummaryBuilderFactory
}

//NewRunCommand ...
func NewRunCommand(app *kingpin.Application, registry *core.Registry, summaryBuilders *core.SummaryBuilderFactory) {
	configuration := &config.Configuration{}
	c := &RunCommand{
		Config:          configuration,
		registry:        registry,
		summaryBuilders: summaryBuilders,
	}
	run := app.Command("run", "Execute performance test thing").Action(c.run)
	run.Arg("file", "Corcel file contains URLs or an ExecutionPlan (see the --plan argument)").Required().StringVar(&configuration.FilePath)
	run.Flag("summary", "Output summary to STDOUT").BoolVar(&configuration.Summary)
	run.Flag("summary-format", "Format for the summary").Default("console").EnumVar(&configuration.SummaryFormat, "console", "json", "yaml")
	run.Flag("iterations", "The number of iterations to run").Short('i').Default("0").IntVar(&configuration.Iterations)
	run.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").Short('d').Default("0s").DurationVar(&configuration.Duration)
	run.Flag("wait-time", "Time to wait between each execution").Default("0s").Short('t').DurationVar(&configuration.WaitTime)
	run.Flag("workers", "The number of workers to execute the requests").Short('w').IntVar(&configuration.Workers)
	run.Flag("random", "Select the url at random for each execution").Short('r').BoolVar(&configuration.Random)
	run.Flag("plan", "Indicate that the corcel file is an ExecutionPlan").BoolVar(&configuration.Plan)
	run.Flag("verbose", "verbosity").Short('v').Action(config.Counter).Bool()
	run.Flag("progress", "Progress reporter").EnumVar(&configuration.Progress, "bar", "logo", "none")
}

func (this *RunCommand) run(c *kingpin.ParseContext) error {
	configuration, err := config.ParseConfiguration(this.Config)
	if err != nil {
		errormanager.Log(err)
		panic("REPLACE ME THIS IS TEMP")
		os.Exit(1)
	}
	logger.ConfigureLogging(configuration)

	//This will not be anything other than a Console Host as we are working with Run command and Server command. In essence being in this method means we are inside the Console Host.
	app := Application{}
	output := app.Execute(configuration)

	reporter := report.CreateHTMLReporter()
	reporter.Generate(output)

	if configuration.Summary {
		summary := output.CreateSummary()
		summaryBuilder := this.summaryBuilders.Get(configuration.SummaryFormat)
		summaryBuilder.Write(summary)
	}
	return nil
}

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}
