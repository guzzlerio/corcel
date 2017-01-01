package logger

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/guzzlerio/corcel/config"
)

var (
	//Log ...
	Log         *log.Logger
	initialised = false
)

//Initialise ...
func Initialise() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	Log = log.New()
	initialised = true
}

//ConfigureLogging ...
func ConfigureLogging(config *config.Configuration) {
	if !initialised {
		panic("You need to logger.Initialise() first")
	}
	Log.Level = config.LogLevel
	//TODO probably have another ticket to support outputting logs to a file
	//Log.Formatter = config.Logging.Formatter
}
