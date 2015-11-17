package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	applicationVersion = "0.1.1-alpha"
	logEnabled         = false
	//Log ...
	Log *log.Logger
	//RandomSource ...
	RandomSource = rand.NewSource(time.Now().UnixNano())
	//Random ...
	Random = rand.New(RandomSource)
	//ErrorMappings ...
	ErrorMappings = map[string]ErrorCode{}
)

func check(err error) {
	if err != nil {
		for mapping, errorCode := range ErrorMappings {
			if strings.Contains(fmt.Sprintf("%v", err), mapping) {
				fmt.Println(errorCode.Message)
				os.Exit(errorCode.Code)
			}
		}
		Log.Fatalf("UNKNOWN ERROR: %v", err)
	}
}

//ConfigureLogging ...
func ConfigureLogging(config *Configuration) {
	Log = log.New()
	Log.Level = config.LogLevel
	if !logEnabled {
		Log.Out = ioutil.Discard
	}
	//TODO probably have another ticket to support outputting logs to a file
	//Log.Formatter = config.Logging.Formatter
}

//ExecuteRequest ...
func ExecuteRequest(client *http.Client, stats *Statistics, request *http.Request) {
	start := time.Now()
	response, responseError := client.Do(request)
	duration := time.Since(start) / time.Millisecond
	check(responseError)

	defer func() {
		err := response.Body.Close()
		if err != nil {
			Log.Printf("Error closing response Body %v", err)
		}
	}()
	responseBytes, _ := httputil.DumpResponse(response, true)
	stats.BytesReceived(int64(len(responseBytes)))
	if response.StatusCode >= 400 && response.StatusCode < 600 {
		responseError = errors.New("5XX Response Code")
	}

	stats.ResponseTime(int64(duration))
	requestBytes, _ := httputil.DumpRequest(request, true)
	stats.BytesSent(int64(len(requestBytes)))
	stats.Request(responseError)
}

//Execute ...
func Execute(config *Configuration, stats *Statistics) {
	var waitGroup sync.WaitGroup

	reader := NewRequestReader(config.FilePath)

	for i := 0; i < config.Workers; i++ {
		waitGroup.Add(1)
		go func() {
			defer func() { //catch or finally
				if err := recover(); err != nil { //catch
					if strings.Contains(fmt.Sprintf("%v", err), "too many open files") {
						Log.Fatalf("Too many workers man!")
					} else {
						Log.Fatalf("UNKNOWN ERROR: %v", err)
					}
				}
			}()
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 50,
				},
			}
			var stream RequestStream

			if config.Random {
				stream = NewRandomRequestStream(reader)
			} else {
				stream = NewSequentialRequestStream(reader)
			}
			if config.Duration > 0 {
				stream = NewTimeBasedRequestStream(stream, config.Duration)
			}
			for stream.HasNext() {
				request, err := stream.Next()
				check(err)
				ExecuteRequest(client, stats, request)

				time.Sleep(config.WaitTime)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}

func GenerateExecutionOutput(file string, stats *Statistics) {
	outputPath, err := filepath.Abs(file)
	check(err)
	output := stats.ExecutionOutput()
	yamlOutput, err := yaml.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

//OutputSummary ...
func OutputSummary(stats *Statistics) {
	output := stats.ExecutionOutput()
	fmt.Println(fmt.Sprintf("Running Time: %v s", output.Summary.RunningTime/1000))
	fmt.Println(fmt.Sprintf("Throughput: %v req/s", int64(output.Summary.Requests.Rate)))
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
	fmt.Println(fmt.Sprintf("Max Response Time: %v ms", output.Summary.ResponseTime.Max))
}

func main() {
	config, err := ParseConfiguration(os.Args[1:])
	check(err)

	configureErrorMappings()
	ConfigureLogging(config)


	absolutePath, err := filepath.Abs(config.FilePath)
	check(err)
	file, err := os.Open(absolutePath)
	defer func() {
		err := file.Close()
		if err != nil {
			Log.Printf("Error closing file %v", err)
		}
	}()
	check(err)

	stats := CreateStatistics()

	Execute(config, stats)

	stats.Stop()

	check(err)
	GenerateExecutionOutput("./output.yml", stats)

	if config.Summary {
		OutputSummary(stats)
	}
}
