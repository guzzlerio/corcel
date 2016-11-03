package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/guzzlerio/corcel/errormanager"
	"gopkg.in/yaml.v2"
)

//CreateFileFromLines ...
func CreateFileFromLines(lines []string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	errormanager.Check(err)
	for _, line := range lines {
		_, err := file.WriteString(fmt.Sprintf("%s\n", line))
		errormanager.Check(err)
	}
	errormanager.Check(file.Sync())
	return file
}

//PathExists ...
func PathExists(value string) bool {
	path, pathErr := filepath.Abs(value)
	errormanager.Check(pathErr)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//UnmarshalYamlFromFile ...
func UnmarshalYamlFromFile(path string, output interface{}) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, output)
	if err != nil {
		panic(err)
	}
}
