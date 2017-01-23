package utils

import (
	"os"
	"path/filepath"
	"strings"
)

//FindFileUp ...
func FindFileUp(fileName string) string {

	wd, _ := os.Getwd()

	var split = strings.Split(wd, string(os.PathSeparator))

	if split[0] == "" {
		split[0] = "/"
	}

	var returnArray = []string{}

	var depth = 1
	for depth <= len(split) {

		candidate := filepath.Join(strings.Join(split[0:depth], string(os.PathSeparator)), fileName)
		returnArray = append([]string{candidate}, returnArray...)
		depth = depth + 1
	}

	for _, item := range returnArray {
		if _, err := os.Stat(item); !os.IsNotExist(err) {
			return item
		}
	}

	return ""

}
