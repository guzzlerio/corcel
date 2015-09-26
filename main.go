package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	logEnabled = false
	Log        *log.Logger
)

func check(err error) {
	if err != nil {
		fmt.Errorf("%v", err)
		panic(err.Error())
	}
}

func configureLogging() {
	//TODO: refine this to work with levels or replace
	//with a package which already handles this
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	prefix := "cns: "
	if logEnabled {
		Log = log.New(os.Stdout, prefix, flags)
	} else {
		//Send all the output to dev null
		Log = log.New(ioutil.Discard, prefix, flags)
	}
}

func main() {
	filePath := kingpin.Flag("file", "Urls file").Short('f').String()
	kingpin.Parse()

	configureLogging()

	absolutePath, err := filepath.Abs(*filePath)
	check(err)
	file, err := os.Open(absolutePath)
	check(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	client := &http.Client{}
	requestAdapter := NewRequestAdapter()
	for scanner.Scan() {
		line := scanner.Text()
		req, err := requestAdapter.Create(line)
		if err != nil {
			panic(err)
		}
		_, _ = client.Do(req)
	}

	summaryPath, err := filepath.Abs("./output.json")
	if err != nil{
		panic(err)
	}
	err = ioutil.WriteFile(summaryPath, []byte(""), 0644)

	if err != nil{
		panic(err)
	}
}
