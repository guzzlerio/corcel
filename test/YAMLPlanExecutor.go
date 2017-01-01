package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/spf13/hugo/utils"
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

	//path := "./corcel"
	exePath, err := filepath.Abs(path)
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
	cmd := exec.Command(exePath, append(append([]string{"run", "--progress", "none"}, args...), file.Name())...)
	output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	logger.Log.Println(fmt.Sprintf("%s", output))
	return output, err
}

//ExecutePlanBuilder ...
func ExecutePlanBuilder(path string, planBuilder *yaml.PlanBuilder) ([]byte, error) {

	file, err := planBuilder.Build()
	if err != nil {
		return []byte{}, err
	}

	//path := "./corcel"
	exePath, err := filepath.Abs(path)
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
	cmd := exec.Command(exePath, append(append([]string{"run", "--progress", "none"}, args...), file.Name())...)
	output, err := cmd.CombinedOutput()
	//fmt.Println(fmt.Sprintf("OUTPUT: %v\nERROR: %v\n", string(output), err))
	logger.Log.Println(fmt.Sprintf("%s", output))
	return output, err
}
