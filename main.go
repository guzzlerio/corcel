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
	"github.com/gosuri/uiprogress"
	"gopkg.in/yaml.v2"
)

var (
	applicationVersion = "0.1.1-alpha"
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
	//TODO probably have another ticket to support outputting logs to a file
	//Log.Formatter = config.Logging.Formatter
}

//ExecuteRequest ...
func ExecuteRequest(client *http.Client, stats *Statistics, request *http.Request) {
	Log.Infof("%s to %s", request.Method, request.URL)
	start := time.Now()
	response, responseError := client.Do(request)
	duration := time.Since(start) / time.Millisecond
	check(responseError)

	defer func() {
		err := response.Body.Close()
		if err != nil {
			Log.Warnf("Error closing response Body %v", err)
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
	uiprogress.Start()
	bar := uiprogress.AddBar(100).AppendCompleted()

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

				bar.Set(stream.Progress())

				time.Sleep(config.WaitTime)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}

//GenerateExecutionOutput ...
func GenerateExecutionOutput(file string, stats *Statistics) {
	outputPath, err := filepath.Abs(file)
	check(err)
	output := stats.ExecutionOutput()
	yamlOutput, err := yaml.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
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
	stats.Start()

	Execute(config, stats)

	stats.Stop()

	check(err)
	GenerateExecutionOutput("./output.yml", stats)

	if config.Summary {
		output := stats.ExecutionOutput()
		consoleWriter := ExecutionOutputWriter{output}
		consoleWriter.Write(os.Stdout)
	}
}
