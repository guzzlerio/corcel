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
	config, err := config.ParseConfiguration(os.Args[1:])
	if err != nil {
		logger.Log.Fatal(err)
	}

	logger.ConfigureLogging(config)

	absolutePath, err := filepath.Abs(config.FilePath)
	check(err)
	file, err := os.Open(absolutePath)
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Log.Printf("Error closing file %v", err)
		}
	}()
	check(err)

	host := cmd.NewConsoleHost(config)
	id, _ := host.Control.Start(config) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	GenerateExecutionOutput("./output.yml", output)

	if config.Summary {
		consoleWriter := processor.ExecutionOutputWriter{
			Output: output,
		}
		consoleWriter.Write(os.Stdout)
	}
}
