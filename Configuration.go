package main

import (
	"fmt"
	"io/ioutil"
	//"log"
	"os"
	"path"
	"time"

	"github.com/imdario/mergo"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Duration time.Duration
	FilePath string
	Random   bool
	Summary  bool
	Workers  int `yaml:"workers"`
	WaitTime time.Duration
}

func ParseConfiguration(args []string) (*Configuration, error) {
	config := defaultConfig()
	cmd, err := cmdConfig(args)
	if err != nil {
		return nil, err
	}

	pwd, err := pwdConfig()
	if err != nil {
		fmt.Println("boom")
		return nil, err
	}
	fmt.Printf("\ndefault:  %+v\n", config)
	fmt.Printf("\n    cmd:  %+v\n", cmd)
	fmt.Printf("\n    pwd:  %+v\n", pwd)
	//usr, err := userDirConfig()
	if err := mergo.Merge(&config, &cmd); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&pwd, &config); err != nil {
		return nil, err
	}
	fmt.Printf("\n config:  %+v\n", pwd)
	return &pwd, err
}

/*
Might need to change the approach here. Issue is if you have a .corcelrc with workers: 3
but you invoke with --workers 1 becuase 1 is the default it will not know to apply the commandline option and will therefore apply the rc file.

One option is to use kingpin.ExpandArgsFromFile

Another option is to use a map, although this wouldn't solve the issue above
*/
func cmdConfig(args []string) (Configuration, error) {
	CommandLine := kingpin.New("corcel", "")
	filePath := CommandLine.Arg("file", "Urls file").Required().ExistingFile()
	summary := CommandLine.Flag("summary", "Output summary to STDOUT").Bool()
	waitTimeArg := CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").String()
	workers := CommandLine.Flag("workers", "The number of workers to execute the requests").Default("1").Int()
	random := CommandLine.Flag("random", "Select the url at random for each execution").Bool()
	durationArg := CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").String()

	fmt.Printf("Parsing cmd args: %+v\n", args)
	cmd, err := CommandLine.Parse(args)
	fmt.Printf("Parsed: %s %s", cmd, err)

	if err != nil {
		fmt.Println(err)
		fmt.Println(cmd)
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
	fmt.Println(contents)
	var config Configuration
	if err := config.Parse(contents); err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func userDirConfig() Configuration {
	return Configuration{}
}

func defaultConfig() Configuration {
	waitTime, _ := time.ParseDuration("0s")
	duration := time.Duration(0)
	return Configuration{
		Duration: duration,
		Random:   false,
		Summary:  false,
		//This seems a little odd, but the merge treats a non-zero value as being set
		Workers:  0,
		WaitTime: waitTime,
	}
}

func (c *Configuration) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	fmt.Println(c, string(data[:]))
	/*
		if c.Hostname == "" {
			return errors.New("Kitchen config: invalid `hostname`")
		}
	*/
	return nil
}

var configFileReader = func(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(fmt.Errorf("not found %s", path))
		return nil, nil
	}
	fmt.Printf("file exists; processing...")
	//TODO actually read the file contents
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s", path)
	}
	return data, nil
}
