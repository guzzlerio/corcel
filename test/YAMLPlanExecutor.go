package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/yaml"
)

//ExecutePlanBuilder ...
func ExecutePlanBuilder(path string, planBuilder *yaml.PlanBuilder) error {

	file, err := planBuilder.BuildAndWrite()
	if err != nil {
		return err
	}

	//path := "./corcel"
	exePath, err := filepath.Abs(path)
	if err != nil {
		return err
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
	return err
}
