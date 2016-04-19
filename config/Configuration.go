package config

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	applicationVersion = "0.1.1-alpha"
	CommandLine        = kingpin.New("corcel", "")
)

//Configuration ...
type Configuration struct {
	Random   bool          `yaml:"random"`
	Summary  bool          `yaml:"summary"`
	LogLevel log.Level     `yaml:"log-level"`
	Workers  int           `yaml:"workers"`
	Duration time.Duration `yaml:"duration"`
	WaitTime time.Duration `yaml:"wait-time"`
	Progress string        `yaml:"progress"`
	Plan     bool          `yaml:"plan"`
	FilePath string
}

func (instance *Configuration) validate() error {
	if err := instance.handleHTTPEndpointForURLFile(); err != nil {
		return err
	}

	if _, err := os.Stat(instance.FilePath); os.IsNotExist(err) {
		return errors.New("required argument 'file' not provided")
	}
	return nil
}

func Usage() {
	CommandLine.Usage(os.Args[1:])
}

//ParseConfiguration ...
func ParseConfiguration(args []string) (*Configuration, error) {
	verbosity = 0
	logLevel = log.FatalLevel
	config := Configuration{}
	defaults := DefaultConfig()
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

	eachConfig := []interface{}{cmd, pwd, usr, &defaults}
	for _, item := range eachConfig {
		if err := mergo.Merge(&config, item); err != nil {
			return nil, err
		}
	}
	setLogLevel(&config, eachConfig)
	log.WithFields(log.Fields{"config": config}).Info("Configuration")

	return &config, err
}

var verbosity int
var logLevel log.Level

func setLogLevel(config *Configuration, each []interface{}) {
	max := log.PanicLevel
	for _, value := range each {
		if value.(*Configuration).LogLevel > max {
			max = value.(*Configuration).LogLevel // found another smaller value, replace previous value in max
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

func cmdConfig(args []string) (*Configuration, error) {
	CommandLine := kingpin.New("corcel", "")
	CommandLine.HelpFlag.Short('h')
	CommandLine.UsageTemplate(kingpin.LongHelpTemplate)

	CommandLine.Version(applicationVersion)

	config := Configuration{}
	CommandLine.Arg("file", "Corcel file contains URLs or an ExecutionPlan (see the --plan argument)").Required().StringVar(&config.FilePath)
	CommandLine.Flag("summary", "Output summary to STDOUT").BoolVar(&config.Summary)
	CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").Default("0s").DurationVar(&config.Duration)
	CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").DurationVar(&config.WaitTime)
	CommandLine.Flag("workers", "The number of workers to execute the requests").IntVar(&config.Workers)
	CommandLine.Flag("random", "Select the url at random for each execution").BoolVar(&config.Random)
	CommandLine.Flag("plan", "Indicate that the corcel file is an ExecutionPlan").BoolVar(&config.Plan)
	CommandLine.Flag("verbose", "verbosity").Short('v').Action(counter).Bool()
	CommandLine.Flag("progress", "Progress reporter").EnumVar(&config.Progress, "bar", "logo", "none")

	_, err := CommandLine.Parse(args)

	if err != nil {
		return nil, err
	}

	if verbosity > 0 {
		config.Progress = "none"
	}

	config.LogLevel = logLevel

	if err = config.handleHTTPEndpointForURLFile(); err != nil {
		return nil, err
	}

	if _, err = os.Stat(config.FilePath); os.IsNotExist(err) {
		return nil, errors.New("required argument 'file' not provided")
	}

	if validationErr := config.validate(); validationErr != nil {
		return nil, validationErr
	}

	return &config, nil
}

func pwdConfig() (*Configuration, error) {
	pwd, _ := os.Getwd()
	// can we get the application name programatically?
	filename := path.Join(pwd, fmt.Sprintf(".%src", "corcel"))

	contents, err := configFileReader(filename)
	if err != nil {
		return nil, err
	}
	var config Configuration
	if err := config.parse(contents); err != nil {
		return nil, err
	}
	return &config, nil
}

func userDirConfig() (*Configuration, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	home, err := homedir.Expand(dir)
	// can we get the application name programatically?
	filename := path.Join(home, fmt.Sprintf(".%src", "corcel"))

	contents, err := configFileReader(filename)
	if err != nil {
		return nil, err
	}
	var config Configuration
	if err := config.parse(contents); err != nil {
		return nil, err
	}
	return &config, nil
}

//DefaultConfig ...
func DefaultConfig() Configuration {
	waitTime := time.Duration(0)
	duration := time.Duration(0)
	return Configuration{
		Duration: duration,
		Plan:     false,
		Random:   false,
		Summary:  false,
		Workers:  1,
		WaitTime: waitTime,
		LogLevel: log.FatalLevel,
		Progress: "logo",
	}
}

func (instance *Configuration) handleHTTPEndpointForURLFile() error {
	u, e := url.ParseRequestURI(instance.FilePath)
	if e == nil && u.Scheme != "" {
		log.Printf("Dowloading URL file from %s ...\n", instance.FilePath)
		file, _ := createTemporaryFile(instance.FilePath)
		out, _ := os.Create(file.Name())
		defer func() {
			check(out.Close())
		}()

		body, e := downloadURLFileFromEndpoint(instance.FilePath)
		if e != nil {
			return fmt.Errorf("unable to download url file from endpoint %s [%s]", instance.FilePath, e)
		}
		defer check(body.Close())
		_, _ = io.Copy(out, body)
		instance.FilePath = file.Name()
	}
	return nil
}

func (instance *Configuration) parse(data []byte) error {
	if err := yaml.Unmarshal(data, instance); err != nil {
		log.Warn("Unable to parse config file")
		return nil
	}
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

var downloadURLFileFromEndpoint = func(endpoint string) (io.ReadCloser, error) {
	resp, e := http.Get(endpoint)
	if e != nil {
		return nil, e
	}
	return resp.Body, nil
}

var createTemporaryFile = func(filePath string) (*os.File, error) {
	hashed := md5.Sum([]byte(filePath))
	return ioutil.TempFile(os.TempDir(), fmt.Sprintf("%x", hashed))
}

func check(err error) {
	if err != nil {
		log.Fatalf("UNKNOWN ERROR: %v", err)
	}
}
