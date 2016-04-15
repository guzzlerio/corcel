package logger

import (
	log "github.com/Sirupsen/logrus"

	"ci.guzzler.io/guzzler/corcel/config"
)

var (
	//Log ...
	Log         *log.Logger
	initialised = false
)

//Initialise ...
func Initialise() {
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
