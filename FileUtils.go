package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//CreateFileFromLines ...
func CreateFileFromLines(lines []string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	check(err)
	for _, line := range lines {
		_, err := file.WriteString(fmt.Sprintf("%s\n", line))
		check(err)
	}
	check(file.Sync())
	return file
}

//PathExists ...
func PathExists(value string) bool {
	path, pathErr := filepath.Abs(value)
	check(pathErr)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//UnmarshalYamlFromFile ...
func UnmarshalYamlFromFile(path string, output interface{}) {
	absPath, err := filepath.Abs(path)
	check(err)
	data, err := ioutil.ReadFile(absPath)
	check(err)
	err = yaml.Unmarshal(data, output)
	check(err)
}
