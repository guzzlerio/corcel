package main

import "time"

type Configuration struct {
	workers int
    waitTime time.Duration
}

func ParseConfiguration() *Configuration {
	return defaultConfig()
}

func defaultConfig() *Configuration {
	return &Configuration{workers: 1}
}
