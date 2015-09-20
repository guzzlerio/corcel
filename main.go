package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"gopkg.in/alecthomas/kingpin.v2"
)

func check(err error) {
	if err != nil {
		fmt.Errorf("%v",err)
		panic(err.Error())
	}
}

func main() {
	filePath := kingpin.Flag("file", "Urls file").Short('f').String()
	kingpin.Parse()

	absolutePath, err := filepath.Abs(*filePath)
	check(err)
	file, err := os.Open(absolutePath)
	check(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	client := &http.Client{}
	for scanner.Scan() {
		line := scanner.Text()
		req, _ := http.NewRequest("GET", line, nil)
		_, _ = client.Do(req)
	}
}
