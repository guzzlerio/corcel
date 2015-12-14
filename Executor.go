package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
	req "ci.guzzler.io/guzzler/corcel/request"
)

var ()

type Executor struct {
	config *config.Configuration
	stats  *Statistics
}

//Execute ...
func (instance *Executor) Execute() {
	var waitGroup sync.WaitGroup

	reader := req.NewRequestReader(instance.config.FilePath)
	bar := NewProgressBar(100, instance.config)

	for i := 0; i < instance.config.Workers; i++ {
		waitGroup.Add(1)
		go func() {
			defer func() { //catch or finally
				if err := recover(); err != nil { //catch
					if strings.Contains(fmt.Sprintf("%v", err), "too many open files") {
						logger.Log.Fatalf("Too many workers man!")
					} else {
						logger.Log.Fatalf("UNKNOWN ERROR: %v", err)
					}
				}
			}()
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 50,
				},
			}
			var stream req.RequestStream

			if instance.config.Random {
				stream = req.NewRandomRequestStream(reader)
			} else {
				stream = req.NewSequentialRequestStream(reader)
			}
			if instance.config.Duration > 0 {
				stream = req.NewTimeBasedRequestStream(stream, instance.config.Duration)
			}
			for stream.HasNext() {
				request, err := stream.Next()
				check(err)
				instance.executeRequest(client, request)

				_ = bar.Set(stream.Progress())

				time.Sleep(instance.config.WaitTime)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}

func (instance *Executor) executeRequest(client *http.Client, request *http.Request) {
	logger.Log.Infof("%s to %s", request.Method, request.URL)
	start := time.Now()
	response, responseError := client.Do(request)
	duration := time.Since(start) / time.Millisecond
	check(responseError)

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Log.Warnf("Error closing response Body %v", err)
		}
	}()
	responseBytes, _ := httputil.DumpResponse(response, true)
	instance.stats.BytesReceived(int64(len(responseBytes)))
	if response.StatusCode >= 400 && response.StatusCode < 600 {
		responseError = errors.New("5XX Response Code")
	}

	instance.stats.ResponseTime(int64(duration))
	requestBytes, _ := httputil.DumpRequest(request, true)
	instance.stats.BytesSent(int64(len(requestBytes)))
	instance.stats.Request(responseError)
}
