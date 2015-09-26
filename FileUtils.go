package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateFileFromLines(lines []string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		file.WriteString(fmt.Sprintf("%s\n", line))
	}
	file.Sync()
	return file
}

func PathExists(value string) bool {
	path, pathErr := filepath.Abs(value)
	if pathErr != nil {
		panic(pathErr)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

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
