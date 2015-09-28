package main

import (
	"fmt"
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"
	"errors"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	logEnabled = false
	Log        *log.Logger
)

func check(err error) {
	if err != nil {
		Log.Panic(err)
	}
}

func ConfigureLogging() {
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

func Execute(file *os.File, stats *Statistics) {
	defer file.Close()
	scanner := bufio.NewScanner(file)

	client := &http.Client{}
	requestAdapter := NewRequestAdapter()

	for scanner.Scan() {
		line := scanner.Text()
		request, err := requestAdapter.Create(line)
		check(err)
		start := time.Now()
		response, err := client.Do(request)
		duration := time.Since(start) / time.Millisecond
		check(err)
		requestBytes, _ := httputil.DumpRequest(request, true)
		responseBytes, _ := httputil.DumpResponse(response, true)

		stats.BytesReceived(int64(len(responseBytes)))
		stats.BytesSent(int64(len(requestBytes)))
		stats.ResponseTime(int64(duration))

		var responseError error = nil
		if (response.StatusCode >= 400 && response.StatusCode < 600) {
			responseError = errors.New("5XX Response Code")
		}
		stats.Request(responseError)
	}
}

func GenerateExecutionOutput(outputPath string, stats *Statistics) {
	output := stats.ExecutionOutput()
	yamlOutput, err := yaml.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

func OutputSummary(stats *Statistics){
	output := stats.ExecutionOutput()
	fmt.Println(fmt.Sprintf("Running Time: %v seconds",output.Summary.RunningTime / 1000))
}

func main() {
	filePath := kingpin.Flag("file", "Urls file").Short('f').String()
	summary := kingpin.Flag("summary", "Output summary to STDOUT").Bool()
	kingpin.Parse()

	ConfigureLogging()

	absolutePath, err := filepath.Abs(*filePath)
	check(err)
	file, err := os.Open(absolutePath)
	check(err)

	stats := CreateStatistics()
	stats.Start()

	Execute(file, stats)

	outputPath, err := filepath.Abs("./output.yml")
	check(err)
	GenerateExecutionOutput(outputPath, stats)

	if *summary {
		OutputSummary(stats)
	}
}
