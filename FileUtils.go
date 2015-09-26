package main

import (
	"io/ioutil"
	"fmt"
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
