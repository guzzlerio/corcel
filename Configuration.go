package main

import "time"

type Configuration struct {
    duration time.Duration
    random bool
    summary bool
    workers int
    waitTime time.Duration
}

func ParseConfiguration() *Configuration {
    return defaultConfig()
}

func defaultConfig() *Configuration {
    waitTime, _ := time.ParseDuration("0s")
    duration := time.Duration(0)
    return &Configuration{
        duration: duration,
        random: false,
        summary: false,
        workers: 1,
        waitTime: waitTime,
    }
}
