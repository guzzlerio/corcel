package main

import "time"

type Configuration struct {
    random bool
    workers int
    waitTime time.Duration
}

func ParseConfiguration() *Configuration {
    return defaultConfig()
}

func defaultConfig() *Configuration {
    waitTime, _ := time.ParseDuration("0s")
    return &Configuration{
        random: false,
        workers: 1,
        waitTime: waitTime,
    }
}
