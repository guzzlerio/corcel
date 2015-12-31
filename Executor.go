package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/processor"
	req "ci.guzzler.io/guzzler/corcel/request"
)

// Executor ...
type Executor struct {
	config *config.Configuration
	stats  *processor.Statistics
	bar    ProgressBar
}

// Execute ...
func (instance *Executor) Execute() {
	instance.stats.Start()
	var waitGroup sync.WaitGroup

	reader := req.NewRequestReader(instance.config.FilePath)

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

				_ = instance.bar.Set(stream.Progress())

				time.Sleep(instance.config.WaitTime)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
	instance.stats.Stop()
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

// Output ...
func (instance *Executor) Output() processor.ExecutionOutput {
	return instance.stats.ExecutionOutput()
}

// ExecutionID ...
type ExecutionID struct {
	value string
}

// String ...
func (id ExecutionID) String() string {
	return fmt.Sprintf("%s", id.value)
}

// NewExecutionID ...
func NewExecutionID() ExecutionID {
	id := randString(16)
	return ExecutionID{id}
}

func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
