package cmd

import (
	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/infrastructure/inproc"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"

	"github.com/rcrowley/go-metrics"
)

//Application ...
type Application struct{}

//Execute ...
func (instance Application) Execute(configuration *config.Configuration) statistics.AggregatorSnapShot {

	metrics.DefaultRegistry.UnregisterAll()

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

	host := NewConsoleHost(configuration, registry)
	id, _ := host.Control.Start(configuration)
	output := host.Control.Stop(id)
	return output
}
