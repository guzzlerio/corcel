package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/imdario/mergo"
	"gopkg.in/alecthomas/kingpin.v2"
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

var (
	applicationVersion = "v0.1.4-alpha"
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

// ServerCommand ...
type ServerCommand struct {
	Port     int
	registry Registry
}

func (cmd *ServerCommand) run(c *kingpin.ParseContext) error {
	// have access to cmd.registry
	//Start HTTP Server
	// construct HTTP Host
	// Start HTTP Host from cmd options
	return nil
}

func configureServerCommand(app *kingpin.Application, registry *Registry) {
	c := &ServerCommand{
		registry: registry,
	}
	server := app.Command("server", "Start HTTP server").Action(c.run)
	server.Flag("port", "Port").Default("54332").IntVar(&c.Port)
}

func main() {
	args := os.Args[1:]

	logger.Initialise()

	//define registry

	configuration := &config.Configuration{}

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(applicationVersion).Author("Andrew Rea").Author("James Allen")
	kingpin.CommandLine.Help = "An example implementation of curl."
	app := kingpin.New("corcel", "")
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.LongHelpTemplate)

	configureServerCommand(app, registry)

	run := app.Command("run", "Execute performance test thing")
	run.Arg("file", "Corcel file contains URLs or an ExecutionPlan (see the --plan argument)").Required().StringVar(&configuration.FilePath)
	run.Flag("summary", "Output summary to STDOUT").BoolVar(&configuration.Summary)
	run.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").Default("0s").DurationVar(&configuration.Duration)
	run.Flag("wait-time", "Time to wait between each execution").Default("0s").DurationVar(&configuration.WaitTime)
	run.Flag("workers", "The number of workers to execute the requests").IntVar(&configuration.Workers)
	run.Flag("random", "Select the url at random for each execution").BoolVar(&configuration.Random)
	run.Flag("plan", "Indicate that the corcel file is an ExecutionPlan").BoolVar(&configuration.Plan)
	run.Flag("verbose", "verbosity").Short('v').Action(config.Counter).Bool()
	run.Flag("progress", "Progress reporter").EnumVar(&configuration.Progress, "bar", "logo", "none")

	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		configuration, err = config.CmdConfig(configuration)
		if err != nil {
			app.Usage(args)
			os.Exit(1)
		}
		//log.SetLevel(logLevel)

		pwd, err := config.PwdConfig()
		if err != nil {
			app.Usage(args)
			os.Exit(1)
		}
		usr, err := config.UserDirConfig()
		if err != nil {
			app.Usage(args)
			os.Exit(1)
		}

		defaults := config.DefaultConfig()
		eachConfig := []interface{}{configuration, pwd, usr, &defaults}
		for _, item := range eachConfig {
			if err := mergo.Merge(configuration, item); err != nil {
				app.Usage(args)
				os.Exit(1)
			}
		}
		config.SetLogLevel(configuration, eachConfig)
		//log.WithFields(log.Fields{"config": config}).Info("Configuration")
		if err != nil {
			app.Usage(args)
			os.Exit(1)
		}

		logger.ConfigureLogging(configuration)

		//TODO: This is not as efficient as it could be for example:
		//Ideally we would only add the HTTP result processor IF an HTTP Action was used
		//Currently every result processor needs to be added.
		//TODO add a ScanForActions .ScanForAssertions .ScanForProcessors .ScanForExtractors .ScanForContexts
		registry := core.CreateRegistry().
			AddActionParser(inproc.YamlDummyActionParser{}).
			AddActionParser(http.YamlHTTPRequestParser{}).
			AddAssertionParser(yaml.ExactAssertionParser{}).
			AddAssertionParser(yaml.NotEqualAssertionParser{}).
			AddAssertionParser(yaml.EmptyAssertionParser{}).
			AddAssertionParser(yaml.NotEmptyAssertionParser{}).
			AddAssertionParser(yaml.GreaterThanAssertionParser{}).
			AddAssertionParser(yaml.GreaterThanOrEqualAssertionParser{}).
			AddAssertionParser(yaml.LessThanAssertionParser{}).
			AddAssertionParser(yaml.LessThanOrEqualAssertionParser{}).
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
