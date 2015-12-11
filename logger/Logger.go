package logger

import (
	log "github.com/Sirupsen/logrus"

	"ci.guzzler.io/guzzler/corcel/config"
)

var (
	//Log ...
	Log *log.Logger
)

//ConfigureLogging ...
func ConfigureLogging(config *config.Configuration) {
	Log = log.New()
	Log.Level = config.LogLevel
	//TODO probably have another ticket to support outputting logs to a file
	//Log.Formatter = config.Logging.Formatter
}

