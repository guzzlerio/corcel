package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guzzlerio/corcel/cmd"
	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/utils"
)

func planDataToFile(platData string) (*os.File, error) {
	file, err := ioutil.TempFile(os.TempDir(), "yamlExecutionPlanForCorcel")
	if err != nil {
		return nil, err
	}
	defer func() {
		utils.CheckErr(file.Close())
	}()

	file.Write([]byte(platData))

	return file, nil
}

//ExecutePlanFromData ...
func ExecutePlanFromData(path string, planData string) ([]byte, error) {
	file, err := planDataToFile(planData)
	if err != nil {
		return []byte{}, err
	}

	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()

	args := []string{"--plan"}
	return executeShell(path, file, args...)
}

//ExecutePlanFromDataForApplication ...
func ExecutePlanFromDataForApplication(path string, planData string, configuration config.Configuration) (statistics.AggregatorSnapShot, error) {
	file, fileErr := planDataToFile(planData)
	if fileErr != nil {
		return statistics.AggregatorSnapShot{}, fileErr
	}

	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()
	configuration.Progress = "none"
	configuration.FilePath = file.Name()
	configuration.Plan = true

	var appConfig, err = config.ParseConfiguration(&configuration)

	if err != nil {
		return statistics.AggregatorSnapShot{}, err
	}

	app := cmd.Application{}
	output := app.Execute(appConfig)

	return output, nil
}

//ExecutePlanBuilder ...
func ExecutePlanBuilder(path string, planBuilder *yaml.PlanBuilder) ([]byte, error) {
	file, err := planBuilder.Build()
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()

	args := []string{"--plan"}
	return executeShell(path, file, args...)
}

//ExecutePlanBuilderForApplication ...
func ExecutePlanBuilderForApplication(path string, planBuilder *yaml.PlanBuilder, configuration config.Configuration) (statistics.AggregatorSnapShot, error) {
	file, fileErr := planBuilder.Build()
	if fileErr != nil {
		return statistics.AggregatorSnapShot{}, fileErr
	}

	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()

	configuration.Progress = "none"
	configuration.FilePath = file.Name()
	configuration.Plan = true

	var appConfig, err = config.ParseConfiguration(&configuration)

	if err != nil {
		return statistics.AggregatorSnapShot{}, err
	}

	app := cmd.Application{}
	output := app.Execute(appConfig)

	return output, nil
}

//ExecuteListForApplication ...
func ExecuteListForApplication(path string, list []string, configuration config.Configuration) (statistics.AggregatorSnapShot, error) {

	file := utils.CreateFileFromLines(list)
	defer func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	}()

	configuration.Progress = "none"
	configuration.FilePath = file.Name()
	var appConfig, err = config.ParseConfiguration(&configuration)

	if err != nil {
		return statistics.AggregatorSnapShot{}, err
	}

	app := cmd.Application{}
	output := app.Execute(appConfig)

	return output, nil
}

//ExecuteList ...
func ExecuteList(path string, list []string, args ...string) ([]byte, error) {

	file := utils.CreateFileFromLines(list)
	defer func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	}()

	return executeShell(path, file, args...)
}

func executeShell(path string, file *os.File, args ...string) ([]byte, error) {
	exePath, exeErr := filepath.Abs(path)
	if exeErr != nil {
		return []byte{}, exeErr
	}
	cmd := exec.Command(exePath, append(append([]string{"run", "--progress", "none"}, args...), file.Name())...)
	output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	if len(output) > 0 {
		logger.Log.Println(fmt.Sprintf("%s", output))
	}
	return output, err
}
