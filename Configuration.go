package main

import (
	"fmt"
	//"log"
	"time"

	"github.com/imdario/mergo"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Configuration struct {
	Duration time.Duration
	FilePath string
	Random   bool
	Summary  bool
	Workers  int
	WaitTime time.Duration
}

func ParseConfiguration(args []string) (*Configuration, error) {
	config, err := cmdConfig(args)
	if err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, pwdConfig()); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, userDirConfig()); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, defaultConfig()); err != nil {
		return nil, err
	}
	//fmt.Printf("\nconfig:  %+v\n", config)
	return &config, err
}

func cmdConfig(args []string) (Configuration, error) {
	CommandLine := kingpin.New("name", "")
	filePath := CommandLine.Arg("file", "Urls file").Required().String()
	summary := CommandLine.Flag("summary", "Output summary to STDOUT").Bool()
	waitTimeArg := CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").String()
	workers := CommandLine.Flag("workers", "The number of workers to execute the requests").Default("1").Int()
	random := CommandLine.Flag("random", "Select the url at random for each execution").Bool()
	durationArg := CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").String()

	cmd, err := CommandLine.Parse(args)

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

func pwdConfig() Configuration {
	return Configuration{}
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
		Workers:  1,
		WaitTime: waitTime,
	}
}
