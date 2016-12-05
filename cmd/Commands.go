package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"gopkg.in/alecthomas/kingpin.v2"
	yamlv2 "gopkg.in/yaml.v2"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/converters"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/report"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
)

type Command interface {
	Run(c *kingpin.ParseContext) error
}

//New...
func New(app *kingpin.Application, registry *core.Registry) {
	// ServerCommand
	sc := &ServerCommand{
		registry: registry,
	}
	server := app.Command("server", "Start HTTP server").Action(sc.run)
	server.Flag("port", "Port").Default("54332").IntVar(&sc.Port)

	// ConvertCommand
	cc := &ConvertCommand{
		registry: registry,
	}
	convert := app.Command("convert", "Convert input log file into a Corcel Plan").Action(cc.Run)
	convert.Arg("input", "Input File").StringVar(&cc.inputFile)
	convert.Arg("output", "Output .plan file").StringVar(&cc.outputFile)
	convert.Flag("base", "Base URL").Short('b').Required().StringVar(&cc.baseUrl)
	convert.Flag("type", "Log File Type").Short('t').Default("iis").EnumVar(&cc.logType, "iis")
	convert.Flag("converter", "Path to custom JavaScript converter").Short('c').ExistingFileVar(&cc.converter)

	// RunCommand
	configuration := &config.Configuration{}
	rc := &RunCommand{
		Config:   configuration,
		registry: registry,
	}
	run := app.Command("run", "Execute performance test thing").Action(rc.run)
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

// ServerCommand ...
type ServerCommand struct {
	Port     int
	registry *core.Registry
}

func (instance *ServerCommand) run(c *kingpin.ParseContext) error {
	// have access to c.registry
	//Start HTTP Server
	// construct HTTP Host
	// Start HTTP Host from cmd options
	fmt.Printf("Would now be starting the HTTP server on %v\n", instance.Port)
	return nil
}

// ConvertCommand ...
type ConvertCommand struct {
	inputFile    string
	outputFile   string
	logType      string
	formatString string
	baseUrl      string
	converter    string
	registry     *core.Registry
}

func (instance *ConvertCommand) Run(c *kingpin.ParseContext) error {
	file, _ := os.Open(instance.inputFile)
	defer file.Close()

	var buf []byte
	if len(instance.converter) > 0 {
		buf, _ = ioutil.ReadFile(instance.converter)
	} else {
		switch instance.logType {
		case "iis":
			buf, _ = ioutil.ReadFile("./converters/parsers/iisParser.js")
		default:
			panic(fmt.Errorf("Unsupported logType: %v", instance.logType))
		}
	}
	u, err := url.Parse(instance.baseUrl)
	if err != nil {
		panic(err)
	}
	converter := converters.NewJsLogConverter(string(buf), u, file)
	plan, err := converter.Convert()
	if err != nil {
		panic(err)
	}
	//TODO the Write function should be extracted.
	planBuilder := yaml.NewPlanBuilder()
	if file, err := planBuilder.Write(plan); err == nil {
		defer func() {
			//TODO only do this if an output file is not specified
			if fileErr := os.Remove(file.Name()); fileErr != nil {
				panic(fileErr)
			}
		}()
		dat, _ := ioutil.ReadFile(file.Name())

		if instance.outputFile != "" {
			_ = ioutil.WriteFile(instance.outputFile, dat, 0644)
			fmt.Printf("Converted the log file %v written to %v\n", instance.inputFile, instance.outputFile)
		} else {
			fmt.Printf("Converted the log file %v\n", instance.inputFile)
			fmt.Println(string(dat))
		}
	}
	return nil
}

// RunCommand ...
type RunCommand struct {
	Config   *config.Configuration
	registry *core.Registry
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

	reporter := report.CreateHTMLReporter()
	reporter.Generate(output)

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
