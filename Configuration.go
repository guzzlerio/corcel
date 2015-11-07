package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Duration time.Duration `yaml:"duration"`
	FilePath string
	Random   bool          `yaml:"random"`
	Summary  bool          `yaml:"summary"`
	Workers  int           `yaml:"workers"`
	WaitTime time.Duration `yaml:"wait-time"`
}

func ParseConfiguration(args []string) (*Configuration, error) {
	config := Configuration{}
	defaults := defaultConfig()
	cmd, err := cmdConfig(args)
	if err != nil {
		return nil, err
	}

	pwd, err := pwdConfig()
	if err != nil {
		return nil, err
	}
	usr, err := userDirConfig()
	if err != nil {
		return nil, err
	}

    Log.Printf(" default: %+v\n", defaults)
    Log.Printf("     cmd: %+v\n", cmd)
    Log.Printf("     pwd: %+v\n", pwd)
    Log.Printf("     usr: %+v\n", usr)

	if err := mergo.Merge(&config, &cmd); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, &pwd); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, &usr); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, &defaults); err != nil {
		return nil, err
	}
    Log.Printf(" config: %+v\n", config)
	return &config, err
}

func cmdConfig(args []string) (Configuration, error) {
	CommandLine := kingpin.New("corcel", "")
	filePath := CommandLine.Arg("file", "Urls file").Required().ExistingFile()
	summary := CommandLine.Flag("summary", "Output summary to STDOUT").Bool()
	waitTimeArg := CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").String()
	workers := CommandLine.Flag("workers", "The number of workers to execute the requests").Int()
	random := CommandLine.Flag("random", "Select the url at random for each execution").Bool()
	durationArg := CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").String()

	_, err := CommandLine.Parse(args)

	if err != nil {
        Log.Println("Unable to parse the kingpin args")
		return Configuration{}, err
	}
	waitTime, err := time.ParseDuration(*waitTimeArg)
	if err != nil {
		return Configuration{}, fmt.Errorf("Cannot parse the value specified for --wait-time: '%v'", *waitTimeArg)
	}
	var duration time.Duration
	//remove this if when issue #17 is completed
	if *durationArg != "" {
		duration, err = time.ParseDuration(*durationArg)
		if err != nil {
			return Configuration{}, fmt.Errorf("Cannot parse the value specified for --duration: '%v'", *durationArg)
		}
	}

	return Configuration{
		Duration: duration,
		FilePath: *filePath,
		Random:   *random,
		Summary:  *summary,
		Workers:  *workers,
		WaitTime: waitTime,
	}, err
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
	if err := config.Parse(contents); err != nil {
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
	if err := config.Parse(contents); err != nil {
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
	}
}

func (c *Configuration) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
        Log.Println("Unable to parse config file")
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
        Log.Println("Config file not found")
		return nil, nil
	}
	Log.Println("file exists; processing...")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	return data, nil
}
