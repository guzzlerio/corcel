package processor

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/logger"
	req "ci.guzzler.io/guzzler/corcel/request"
)

//ExecutionBranch ...
type ExecutionBranch interface {
	Execute() error
	Output() ExecutionOutput
}

// Executor ...
type Executor struct {
	config *config.Configuration
	stats  *Statistics
	bar    ProgressBar
}

// Execute ...
func (instance *Executor) Execute() error {
	instance.stats.Start()
	var waitGroup sync.WaitGroup

	reader := req.NewRequestReader(instance.config.FilePath)

	for i := 0; i < instance.config.Workers; i++ {
		waitGroup.Add(1)
		go func() {
			defer func() { //catch or finally
				if err := recover(); err != nil { //catch
					errormanager.Log(err)
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
				request, _ := stream.Next()
				instance.executeRequest(client, request)

				_ = instance.bar.Set(stream.Progress())

				time.Sleep(instance.config.WaitTime)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
	instance.stats.Stop()
	return nil
}

func (instance *Executor) executeRequest(client *http.Client, request *http.Request) {
	logger.Log.Infof("%s to %s", request.Method, request.URL)
	start := time.Now()
	response, responseError := client.Do(request)
	duration := time.Since(start) / time.Millisecond
	if responseError != nil {
		errormanager.Log(responseError)
		return
	}

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
func (instance *Executor) Output() ExecutionOutput {
	return instance.stats.ExecutionOutput()
}
