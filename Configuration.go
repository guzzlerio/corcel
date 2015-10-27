package main

import (
    "time"

    "github.com/imdario/mergo"
)

type Configuration struct {
    duration time.Duration
    random bool
    summary bool
    workers int
    waitTime time.Duration
}

func ParseConfiguration() *Configuration {
    result := defaultConfig()
    mergo.Merge(result, pwdConfig)
    mergo.Merge(result, userDirConfig)
    return result
}

func pwdConfig() *Configuration {
    return &Configuration{}
}

func userDirConfig() *Configuration {
    return &Configuration{}
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
