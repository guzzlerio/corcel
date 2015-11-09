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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	logEnabled   = false
	Log          *log.Logger
	RandomSource = rand.NewSource(time.Now().UnixNano())
	Random       = rand.New(RandomSource)
)

func check(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}

func ConfigureLogging(config *Configuration) {
	Log = log.New()
	Log.Level = config.LogLevel
	if logEnabled {
		Log.Out = ioutil.Discard
	}
	//TODO probably have another ticket to support outputting logs to a file
	//Log.Formatter = config.Logging.Formatter
}

func ExecuteRequest(client *http.Client, stats *Statistics, request *http.Request) {
	start := time.Now()
	response, responseError := client.Do(request)
	duration := time.Since(start) / time.Millisecond
	if responseError == nil {
		defer response.Body.Close()
		responseBytes, _ := httputil.DumpResponse(response, true)
		stats.BytesReceived(int64(len(responseBytes)))
		if response.StatusCode >= 400 && response.StatusCode < 600 {
			responseError = errors.New("5XX Response Code")
		}
	} else {
		Log.Panicln(fmt.Sprintf("Error: %v", responseError))
	}

	stats.ResponseTime(int64(duration))
	requestBytes, _ := httputil.DumpRequest(request, true)
	stats.BytesSent(int64(len(requestBytes)))
	stats.Request(responseError)
}

func Execute(file *os.File, stats *Statistics, waitTime time.Duration, workers int, random bool, duration time.Duration) {
	stats.Start()
	defer stats.Stop()
	defer file.Close()
	var waitGroup sync.WaitGroup

	reader := NewRequestReader(file.Name())

	for i := 0; i < workers; i++ {
		waitGroup.Add(1)
		go func() {
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 50,
				},
			}
			var stream RequestStream

			if random {
				stream = NewRandomRequestStream(reader)
			} else {
				stream = NewSequentialRequestStream(reader)
			}
			if duration > 0 {
				stream = NewTimeBasedRequestStream(stream, duration)
			}
			for stream.HasNext() {
				request, err := stream.Next()
				if err != nil {
					panic(err)
				}
				ExecuteRequest(client, stats, request)

				time.Sleep(waitTime)
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
	if err != nil {
		log.Fatal(err)
	}

	ConfigureLogging(config)

	absolutePath, err := filepath.Abs(config.FilePath)
	check(err)
	file, err := os.Open(absolutePath)
	defer file.Close()
	check(err)

	stats := CreateStatistics()

	Execute(file, stats, config.WaitTime, config.Workers, config.Random, config.Duration)

	GenerateExecutionOutput("output.yml", stats)

	if config.Summary {
		OutputSummary(stats)
	}
}
