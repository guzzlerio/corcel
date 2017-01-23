package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/dustin/go-humanize"
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

	//This will not be anything other than a Console Host as we are working with Run command and Server command. In essence being in this method means we are inside the Console Host.
	app := Application{}
	output := app.Execute(configuration)

	reporter := report.CreateHTMLReporter()
	reporter.Generate(output)

	if configuration.Summary {
		outputSummary(output)
	}
	return nil
}

func outputSummary(snapshot statistics.AggregatorSnapShot) {
	summary := statistics.CreateSummary(snapshot)

	top(os.Stdout)
	line(os.Stdout, "Running Time", summary.RunningTime)
	line(os.Stdout, "Throughput", fmt.Sprintf("%-.0f req/s", summary.Throughput))
	line(os.Stdout, "Total Requests", fmt.Sprintf("%-.0f", summary.TotalRequests))
	line(os.Stdout, "Number of Errors", fmt.Sprintf("%-.0f", summary.TotalErrors))
	line(os.Stdout, "Availability", fmt.Sprintf("%-.4f%%", summary.Availability))
	line(os.Stdout, "Bytes Sent", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.TotalSent))))
	line(os.Stdout, "Bytes Received", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.TotalReceived))))
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
}

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}
