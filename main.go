package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

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
		if response.StatusCode >= 400 && response.StatusCode < 600 {
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

func OutputSummary(stats *Statistics) {
	output := stats.ExecutionOutput()
	fmt.Println(fmt.Sprintf("Running Time: %v s", output.Summary.RunningTime/1000))
	fmt.Println(fmt.Sprintf("Total Requests: %v", output.Summary.Requests.Total))
	fmt.Println(fmt.Sprintf("Number of Errors: %v", output.Summary.Requests.Errors))
	fmt.Println(fmt.Sprintf("Availability: %v%%", output.Summary.Requests.Availability*100))
	fmt.Println(fmt.Sprintf("Bytes Sent: %v", output.Summary.Bytes.Sent.Sum))
	fmt.Println(fmt.Sprintf("Bytes Received: %v", output.Summary.Bytes.Received.Sum))
	if output.Summary.ResponseTime.Mean > 0 {
		fmt.Println(fmt.Sprintf("Mean Response Time: %.4v ms", output.Summary.ResponseTime.Mean))
	} else {
		fmt.Println(fmt.Sprintf("Mean Response Time: %v ms", output.Summary.ResponseTime.Mean))
	}

	fmt.Println(fmt.Sprintf("Min Response Time: %v ms", output.Summary.ResponseTime.Min))
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

	stats.Stop()

	outputPath, err := filepath.Abs("./output.yml")
	check(err)
	GenerateExecutionOutput(outputPath, stats)

	if *summary {
		OutputSummary(stats)
	}
}
