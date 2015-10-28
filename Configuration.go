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
	Random   bool
	Summary  bool
	Workers  int64
	WaitTime time.Duration
}

func ParseConfiguration(args []string) *Configuration {
	config := cmdConfig(args)
	mergo.Merge(&config, pwdConfig())
	mergo.Merge(&config, userDirConfig())
	mergo.Merge(&config, defaultConfig())
	fmt.Printf("\nconfig:  %+v\n", config)
	return &config
}

func cmdConfig(args []string) Configuration {
	fmt.Println(args)
	CommandLine := kingpin.New("name", "")
	//filePath := CommandLine.Arg("file", "Urls file").Required().String()
	//summary := CommandLine.Flag("summary", "Output summary to STDOUT").Bool()
	waitTimeArg := CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").String()
	workers := CommandLine.Flag("workers", "The number of workers to execute the requests").Default("1").Int64()
	//random := CommandLine.Flag("random", "Select the url at random for each execution").Bool()
	durationArg := CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").String()

	cmd, err := CommandLine.Parse(args)

	if err != nil {
		fmt.Println(err)
		fmt.Println(cmd)
		panic(err)
	}
	waitTime, _ := time.ParseDuration(*waitTimeArg)
	duration, _ := time.ParseDuration(*durationArg)

	return Configuration{
		Duration: duration,
		//Random:   random,
		//Summary:  summary,
		Workers:  *workers,
		WaitTime: waitTime,
	}
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
		Workers:  int64(1),
		WaitTime: waitTime,
	}
}
