package main

import (
	"fmt"
	//"log"
	"time"

	"github.com/imdario/mergo"
)

type Configuration struct {
	Duration time.Duration
	Random   bool
	Summary  bool
	Workers  int64
	WaitTime time.Duration
}

func ParseConfiguration() *Configuration {
	config := cmdConfig()
	mergo.Merge(&config, pwdConfig())
	mergo.Merge(&config, userDirConfig())
	mergo.Merge(&config, defaultConfig())
	fmt.Printf("\nconfig:  %+v\n", config)
	return &config
}

func cmdConfig() Configuration {
	return Configuration{}
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
