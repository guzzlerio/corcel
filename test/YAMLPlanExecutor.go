package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guzzlerio/corcel/cmd"
	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
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

func stringInSlice(value string, slice []string) bool {
	for _, sliceValue := range slice {
		if value == sliceValue {
			return true
		}
	}
	return false
}

func ensureSummaryInArgs(args []string) []string {
	if !stringInSlice("--summary", args) {
		args = append([]string{"--summary"}, args...)
	}

	if !stringInSlice("--summary-format", args) {
		args = append([]string{"--summary-format", "json"}, args...)
	} else {
		panic("For tests only this method sets the --summary-format")
	}
	return args
}

//ExecutePlanFromData ...
func ExecutePlanFromData(planData string, args ...string) (core.ExecutionSummary, error) {

	file, err := planDataToFile(planData)
	if err != nil {
		return core.ExecutionSummary{}, err
	}
	PrintPlan(file)

	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()

	args = ensureSummaryInArgs(args)
	args = append([]string{"--plan"}, args...)
	output, err := executeShell(utils.FindFileUp("corcel"), file, args...)
	if err != nil {
		return core.ExecutionSummary{}, err
	}

	var executionSummary core.ExecutionSummary

	err = json.Unmarshal(output, &executionSummary)
	if err != nil {
		return core.ExecutionSummary{}, err
	}

	return executionSummary, nil
}

//ExecutePlanFromDataForApplication ...
func ExecutePlanFromDataForApplication(planData string) (core.ExecutionSummary, error) {
	var configuration = config.Configuration{}
	file, fileErr := planDataToFile(planData)
	if fileErr != nil {
		return core.ExecutionSummary{}, fileErr
	}
	PrintPlan(file)

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
		return core.ExecutionSummary{}, err
	}

	app := cmd.Application{}
	output := app.Execute(appConfig)
	var summary = statistics.CreateSummary(output)

	return summary, nil
}

//ExecutePlanBuilder ...
func ExecutePlanBuilder(planBuilder *yaml.PlanBuilder) ([]byte, error) {
	file, err := planBuilder.BuildAndSave()
	if err != nil {
		return []byte{}, err
	}
	PrintPlan(file)
	defer func() {
		fileErr := os.Remove(file.Name())
		if fileErr != nil {
			panic(fileErr)
		}
	}()

	args := []string{"--plan"}
	return executeShell(utils.FindFileUp("corcel"), file, args...)
}

//ExecutePlanBuilderForApplication ...
func ExecutePlanBuilderForApplication(planBuilder *yaml.PlanBuilder) (core.ExecutionSummary, error) {
	var configuration = config.Configuration{}
	file, fileErr := planBuilder.BuildAndSave()

	if fileErr != nil {
		return core.ExecutionSummary{}, fileErr
	}

	PrintPlan(file)

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
		return core.ExecutionSummary{}, err
	}

	app := cmd.Application{}
	output := app.Execute(appConfig)
	var summary = statistics.CreateSummary(output)

	return summary, nil
}

func PrintPlan(file *os.File) {
	if os.Getenv("CORCEL_PRINT_PLAN") == "1" {
		data, fileReadErr := ioutil.ReadFile(file.Name())
		if fileReadErr != nil {
			panic(fileReadErr)
		}
		fmt.Println(string(data))
	}
}

//ExecuteListForApplication ...
func ExecuteListForApplication(list []string, configuration config.Configuration) (core.ExecutionSummary, error) {

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
		return core.ExecutionSummary{}, err
	}

	app := cmd.Application{}
	output := app.Execute(appConfig)
	var summary = statistics.CreateSummary(output)

	return summary, nil
}

//ExecuteList ...
func ExecuteList(list []string, args ...string) (core.ExecutionSummary, error) {

	path := utils.FindFileUp("corcel")

	file := utils.CreateFileFromLines(list)
	defer func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	}()

	args = ensureSummaryInArgs(args)
	output, err := executeShell(path, file, args...)
	if err != nil {
		return core.ExecutionSummary{}, errors.New(string(output))
	}
	var executionSummary core.ExecutionSummary

	err = json.Unmarshal(output, &executionSummary)
	if err != nil {
		return core.ExecutionSummary{}, err
	}
	return executionSummary, nil
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
