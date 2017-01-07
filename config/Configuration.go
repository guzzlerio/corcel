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
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	yamlFormat "github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/alecthomas/kingpin.v2"
)

//Configuration ...
type Configuration struct {
	Iterations int           `json:"iterations"`
	Random     bool          `json:"random"`
	Summary    bool          `json:"summary"`
	LogLevel   log.Level     `json:"log-level"`
	Workers    int           `json:"workers"`
	Duration   time.Duration `json:"duration"`
	WaitTime   time.Duration `json:"wait-time"`
	Progress   string        `json:"progress"`
	Plan       bool          `json:"plan"`
	FilePath   string
}

//WithDuration converts a string duration into a time value and adds it to the configuration
func (instance Configuration) WithDuration(duration string) Configuration {
	value, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	instance.Duration = value
	return instance
}

//WithWaitTime ...
func (instance Configuration) WithWaitTime(waitTime string) Configuration {
	value, err := time.ParseDuration(waitTime)
	if err != nil {
		panic(err)
	}
	instance.WaitTime = value
	return instance
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

var verbosity int
var logLevel = log.FatalLevel

//ParseConfiguration ...
func ParseConfiguration(cfg *Configuration) (*Configuration, error) {
	configuration, err := CmdConfig(cfg)
	if err != nil {
		return nil, err
	}
	//log.SetLevel(logLevel)

	pwd, err := PwdConfig()
	if err != nil {
		return nil, err
	}
	usr, err := UserDirConfig()
	if err != nil {
		return nil, err
	}

	defaults := DefaultConfig()
	eachConfig := []interface{}{configuration, pwd, usr, &defaults}
	for _, item := range eachConfig {
		if err := mergo.Merge(configuration, item); err != nil {
			return nil, err
		}
	}
	SetLogLevel(configuration, eachConfig)
	//log.WithFields(log.Fields{"config": config}).Info("Configuration")

	if _, err = filepath.Abs(configuration.FilePath); err != nil {
		return nil, err
	}

	return configuration, nil
}

//SetLogLevel ...
func SetLogLevel(config *Configuration, each []interface{}) {
	max := log.PanicLevel
	for _, value := range each {
		if value.(*Configuration).LogLevel > max {
			max = value.(*Configuration).LogLevel // found another smaller value, replace previous value in max
		}
	}
	config.LogLevel = max
}

//Counter ...
func Counter(c *kingpin.ParseContext) error {
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

//CmdConfig ...
func CmdConfig(config *Configuration) (*Configuration, error) {

	if verbosity > 0 {
		config.Progress = "none"
	}

	config.LogLevel = logLevel

	if err := config.handleHTTPEndpointForURLFile(); err != nil {
		return nil, err
	}

	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		return nil, errors.New("required argument 'file' not provided")
	}

	if validationErr := config.validate(); validationErr != nil {
		return nil, validationErr
	}

	return config, nil
}

//PwdConfig ...
func PwdConfig() (*Configuration, error) {
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

//UserDirConfig ...
func UserDirConfig() (*Configuration, error) {
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
	if err := yamlFormat.Unmarshal(data, instance); err != nil {
		log.Warn("Unable to parse config file")
		return nil
	}
	return nil
}

var configFileReader = func(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{"path": path}).Debug("Config file not found")
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
