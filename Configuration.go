package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

//Configuration ...
type Configuration struct {
	Duration time.Duration `yaml:"duration"`
	FilePath string
	Random   bool          `yaml:"random"`
	Summary  bool          `yaml:"summary"`
	Workers  int           `yaml:"workers"`
	WaitTime time.Duration `yaml:"wait-time"`
	LogLevel log.Level     `yaml:"log-level"`
}

func ParseConfiguration(args []string) (*Configuration, error) {
    verbosity = 0
	logLevel = log.FatalLevel
	config := Configuration{}
	defaults := defaultConfig()
	cmd, err := cmdConfig(args)
	if err != nil {
		return nil, err
	}
	log.SetLevel(logLevel)

	pwd, err := pwdConfig()
	if err != nil {
		return nil, err
	}
	usr, err := userDirConfig()
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{"default": defaults, "cmd": cmd, "pwd": pwd, "usr": usr}).Debug("Each config")

	if err := mergo.Merge(&config, &cmd); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"config": config}).Debug("Applied cmd args")
	if err := mergo.Merge(&config, &pwd); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"config": config}).Debug("Applied pwd config file")
	if err := mergo.Merge(&config, &usr); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"config": config}).Debug("Applied usr config file")
	if err := mergo.Merge(&config, &defaults); err != nil {
		return nil, err
	}
	eachConfig := []log.Level{cmd.LogLevel, pwd.LogLevel, usr.LogLevel, defaults.LogLevel}
	setLogLevel(&config, eachConfig)
	log.WithFields(log.Fields{"config": config}).Info("Configuration")
	return &config, err
}

var verbosity int
var logLevel log.Level

func setLogLevel(config *Configuration, each []log.Level) {
	max := log.PanicLevel
	for _, value := range each {
		if value > max {
			max = value // found another smaller value, replace previous value in max
		}
	}
	config.LogLevel = max
}

func counter(c *kingpin.ParseContext) error {
	verbosity++
	switch verbosity {
	case 1:
		logLevel = log.WarnLevel
	case 2:
		logLevel = log.InfoLevel
	case 3:
		logLevel = log.DebugLevel
	}
	return nil
}

func cmdConfig(args []string) (Configuration, error) {
	config := Configuration{}
	CommandLine := kingpin.New("corcel", "")

    CommandLine.Version(applicationVersion)

	CommandLine.Arg("file", "Urls file").Required().ExistingFileVar(&config.FilePath)
	CommandLine.Flag("summary", "Output summary to STDOUT").BoolVar(&config.Summary)
	CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").Default("0s").DurationVar(&config.Duration)
	CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").DurationVar(&config.WaitTime)
	CommandLine.Flag("workers", "The number of workers to execute the requests").IntVar(&config.Workers)
	CommandLine.Flag("random", "Select the url at random for each execution").BoolVar(&config.Random)
	CommandLine.Flag("verbose", "verbosity").Short('v').Action(counter).Bool()

	//log.WithFields(log.Fields{ "args": args }).Debug("Parsing cmd args")
	_, err := CommandLine.Parse(args)

	if err != nil {
		log.Error(err)
		return Configuration{}, err
	}
	config.LogLevel = logLevel

	return config, err
}

func pwdConfig() (Configuration, error) {
	pwd, _ := os.Getwd()
	// can we get the application name programatically?
	filename := path.Join(pwd, fmt.Sprintf(".%src", "corcel"))

	contents, err := configFileReader(filename)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	if err := config.parse(contents); err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func userDirConfig() (Configuration, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return Configuration{}, err
	}
	home, err := homedir.Expand(dir)
	// can we get the application name programatically?
	filename := path.Join(home, fmt.Sprintf(".%src", "corcel"))

	contents, err := configFileReader(filename)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	if err := config.parse(contents); err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func defaultConfig() Configuration {
	waitTime := time.Duration(0)
	duration := time.Duration(0)
	return Configuration{
		Duration: duration,
		Random:   false,
		Summary:  false,
		Workers:  1,
		WaitTime: waitTime,
		LogLevel: log.FatalLevel,
	}
}

func (c *Configuration) parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		log.Warn("Unable to parse config file")
		return nil
	}
	/*
		if c.Hostname == "" {
			return errors.New("Kitchen config: invalid `hostname`")
		}
	*/
	return nil
}

var configFileReader = func(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{"path": path}).Warn("Config file not found")
		return nil, nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithFields(log.Fields{"path": path}).Warn("Unable to read config file")
		return nil, nil
	}
	return data, nil
}
