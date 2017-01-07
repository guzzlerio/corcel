package cmd

import (
	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/statistics"
)

//Application ...
type Application struct {
	registry *core.Registry
}

//Execute ...
func (instance Application) Execute(configuration *config.Configuration) statistics.AggregatorSnapShot {
	host := NewConsoleHost(configuration, *instance.registry)
	id, _ := host.Control.Start(configuration)
	output := host.Control.Stop(id)
	return output
}
