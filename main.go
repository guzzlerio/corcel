package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"sync"
	"time"
    "strings"

	"gopkg.in/yaml.v2"
)

var (
	logEnabled   = false
	Log          *log.Logger
	RandomSource = rand.NewSource(time.Now().UnixNano())
	Random       = rand.New(RandomSource)
    ErrorMappings = map[string]ErrorCode{}
)

func check(err error) {
	if err != nil {
        for mapping, errorCode := range ErrorMappings{
            if strings.Contains(fmt.Sprintf("%v",err), mapping){
                fmt.Println(errorCode.Message)
                os.Exit(errorCode.Code)
            }
        }
        Log.Fatalf("UNKNOWN ERROR: %v",err)
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

func ExecuteRequest(client *http.Client, stats *Statistics, request *http.Request) {
	start := time.Now()
	response, responseError := client.Do(request)
	duration := time.Since(start) / time.Millisecond
    check(responseError)

    defer response.Body.Close()
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

func Execute(file *os.File, stats *Statistics, waitTime time.Duration, workers int, random bool, duration time.Duration) {
	defer file.Close()
	var waitGroup sync.WaitGroup

	reader := NewRequestReader(file.Name())

	for i := 0; i < workers; i++ {
		waitGroup.Add(1)
		go func() {
            defer func() { //catch or finally
                if err := recover(); err != nil { //catch
                    if strings.Contains(fmt.Sprintf("%v",err),"too many open files"){
                        Log.Fatalf("Too many workers man!")
                    }else{
                        Log.Fatalf("UNKNOWN ERROR: %v",err)
                    }
                }
            }()
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
                check(err)
				ExecuteRequest(client, stats, request)

				time.Sleep(waitTime)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
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
    ConfigureErrorMappings()
	ConfigureLogging()

	config, err := ParseConfiguration(os.Args[1:])
	check(err)


	absolutePath, err := filepath.Abs(config.FilePath)
	check(err)
	file, err := os.Open(absolutePath)
	defer file.Close()
	check(err)

	stats := CreateStatistics()
	stats.Start()

	Execute(file, stats, config.WaitTime, config.Workers, config.Random, config.Duration)

	stats.Stop()

	outputPath, err := filepath.Abs("./output.yml")
	check(err)
	GenerateExecutionOutput(outputPath, stats)

	if config.Summary {
		OutputSummary(stats)
	}
}
