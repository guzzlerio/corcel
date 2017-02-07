package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/guzzlerio/corcel/cmd"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/infrastructure/inproc"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/json"
	"github.com/guzzlerio/corcel/serialisation/yaml"
)

var (
	//Version ...
	Version = "EMPTY"
	//BuildTime ...
	BuildTime = "EMPTY"
	//CommitHash ...
	CommitHash = "EMPTY"

	cpuprofile = ""

	memprofile = ""

	panicNotRecover = ""
)

func main() {

	defer errormanager.HandlePanic()

	if panicNotRecover != "" {
		errormanager.PanicNotRecover()
	}

	if cpuprofile != "" {
		f, err := os.Create("./corcel.prof")
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	logger.Initialise()

	//TODO: This is not as efficient as it could be for example:
	//Ideally we would only add the HTTP result processor IF an HTTP Action was used
	//Currently every result processor needs to be added.
	//TODO add a ScanForActions .ScanForAssertions .ScanForProcessors .ScanForExtractors .ScanForContexts
	registry := core.CreateRegistry().
		AddActionParser(inproc.YamlDummyActionParser{}).
		AddActionParser(inproc.YamlIPanicActionParser{}).
		AddActionParser(http.YamlHTTPRequestParser{}).
		AddAssertionParser(yaml.ExactAssertionParser{}).
		AddAssertionParser(yaml.NotEqualAssertionParser{}).
		AddAssertionParser(yaml.EmptyAssertionParser{}).
		AddAssertionParser(yaml.NotEmptyAssertionParser{}).
		AddAssertionParser(yaml.GreaterThanAssertionParser{}).
		AddAssertionParser(yaml.GreaterThanOrEqualAssertionParser{}).
		AddAssertionParser(yaml.LessThanAssertionParser{}).
		AddAssertionParser(yaml.LessThanOrEqualAssertionParser{}).
		AddResultProcessor(http.NewExecutionResultProcessor()).
		AddResultProcessor(inproc.NewGeneralExecutionResultProcessor()).
		AddExtractorParser(yaml.RegexExtractorParser{}).
		AddExtractorParser(yaml.XPathExtractorParser{}).
		AddExtractorParser(yaml.JSONPathExtractorParser{}).
		AddExtractorParser(yaml.KeyValueExtractorParser{})

	summaryBuilders := core.NewSummaryBuilderFactory().
		AddBuilder("json", &json.JSONSummaryBuilder{Writer: os.Stdout}).
		AddBuilder("yaml", &yaml.YAMLSummaryBuilder{Writer: os.Stdout}).
		AddBuilder("console", cmd.NewConsoleSummaryBuilder(os.Stdout))

	//kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.Help = "An example implementation of curl."

	versionString := fmt.Sprintf("Version %s, Build Time: %s, Hash: %s", Version, BuildTime, CommitHash)

	app := kingpin.New("corcel", "").Version(versionString).Author("Andrew Rea").Author("James Allen")
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.LongHelpTemplate)

	cmd.NewRunCommand(app, &registry, summaryBuilders)
	cmd.NewServerCommand(app, &registry)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if memprofile != "" {
		f, err := os.Create("./corcel.mprof")
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		err = f.Close()

		if err != nil {
			log.Fatal(err)
		}
		return
	}

}
