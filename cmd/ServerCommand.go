package cmd

import (
	"fmt"

	"github.com/guzzlerio/corcel/core"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
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
