package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
)

var (
	//ErrorMappings ...
	ErrorMappings = map[string]ErrorCode{}
)

func check(err error) {
	if err != nil {
		for mapping, errorCode := range ErrorMappings {
			if strings.Contains(fmt.Sprintf("%v", err), mapping) {
				fmt.Println(errorCode.Message)
				os.Exit(errorCode.Code)
			}
		}
		logger.Log.Fatalf("UNKNOWN ERROR: %v", err)
	}
}

//GenerateExecutionOutput ...
func GenerateExecutionOutput(file string, output ExecutionOutput) {
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

	configureErrorMappings()
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

	host := NewConsoleHost(config)
	id, _ := host.Control.Start(config) //will this block?
	output := host.Control.Stop(id)

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	GenerateExecutionOutput("./output.yml", output)

	if config.Summary {
		consoleWriter := ExecutionOutputWriter{output}
		consoleWriter.Write(os.Stdout)
	}
}
