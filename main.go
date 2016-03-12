package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"ci.guzzler.io/guzzler/corcel/cmd"
	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/processor"
)

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}

//GenerateExecutionOutput ...
func GenerateExecutionOutput(file string, output processor.ExecutionOutput) {
	outputPath, err := filepath.Abs(file)
	check(err)
	yamlOutput, err := yaml.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

func main() {
	logger.Initialise()
	configuration, err := config.ParseConfiguration(os.Args[1:])
	if err != nil {
		config.Usage()
		os.Exit(1)
	}

	logger.ConfigureLogging(configuration)

	_, err = filepath.Abs(configuration.FilePath)
	check(err)

	host := cmd.NewConsoleHost(configuration)
	id, _ := host.Control.Start(configuration) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	GenerateExecutionOutput("./output.yml", output)

	if configuration.Summary {
		consoleWriter := processor.ExecutionOutputWriter{
			Output: output,
		}
		consoleWriter.Write(os.Stdout)
	}
}
