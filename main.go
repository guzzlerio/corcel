package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/guzzlerio/corcel/cmd"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/infrastructure/inproc"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/yaml"
)

var (
	//Version is the application version - set with the ldflags
	Version = ""
)

func main() {
	logger.Initialise()

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
		AddResultProcessor(http.NewExecutionResultProcessor()).
		AddResultProcessor(inproc.NewGeneralExecutionResultProcessor()).
		AddExtractorParser(yaml.RegexExtractorParser{}).
		AddExtractorParser(yaml.XPathExtractorParser{}).
		AddExtractorParser(yaml.JSONPathExtractorParser{})

	//kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.Help = "An example implementation of curl."
	app := kingpin.New("corcel", "").Version(Version).Author("Andrew Rea").Author("James Allen")
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.LongHelpTemplate)

	cmd.NewRunCommand(app, &registry)
	cmd.NewServerCommand(app, &registry)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
