package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/logger"
)

//Action ...
type HTTPAction struct {
	client    *http.Client
	extractor defaultsExtractor
	State     HttpActionState
}

func CreateAction() HTTPAction {
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   0,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 10,
	}
	var instance = HTTPAction{
		client:    &http.Client{Transport: tr},
		extractor: NewDefaultsExtractor(),
		State:     HttpActionState{},
	}
	return instance
}

//Execute ...
func (instance HTTPAction) Execute(ctx context.Context, executionContext core.ExecutionContext) core.ExecutionResult {
	result := core.ExecutionResult{}

	select {
	case <-ctx.Done():
		return result

	default:

		var defaults = instance.extractor.Extract(executionContext)

		if instance.State.Body != "" && instance.State.Body[0] == '@' {
			contents, err := ioutil.ReadFile(instance.State.Body[1:])
			if err != nil {
				result[core.ErrorUrn.String()] = err
				return result
			}
			instance.State.Body = string(contents)
		}

		var requestURL = defaults.URL
		if instance.State.URL != "" {
			requestURL = instance.State.URL
		}

		var method = defaults.Method
		if instance.State.Method != "" {
			method = instance.State.Method
		}

		var headers = defaults.Headers

		var body = defaults.Body
		if instance.State.Body != "" {
			body = instance.State.Body
		}

		for k := range instance.State.Headers {
			headers.Set(k, instance.State.Headers.Get(k))
		}

		for k, v := range executionContext {
			token := k
			switch value := v.(type) {
			case string:
				for hK := range headers {
					replacement := strings.Replace(headers.Get(hK), token, value, -1)
					headers.Set(hK, replacement)
				}
				requestURL = strings.Replace(requestURL, token, value, -1)
				body = strings.Replace(body, token, value, -1)
			}
		}

		requestBody := bytes.NewBuffer([]byte(body))
		req, err := http.NewRequest(method, requestURL, requestBody)
		req = req.WithContext(ctx)
		//req.Cancel = cancellation
		//This should be a configuration item.  It allows the client to work
		//in a way similar to a server which does not support HTTP KeepAlive
		//After each request the client channel is closed.  When set to true
		//the performance overhead is large in terms of Network IO throughput

		req.Close = true

		if err != nil {
			result[core.ErrorUrn.String()] = err
			return result
		}

		req.Header = headers

		response, err := instance.client.Do(req)
		if err != nil {
			result[core.ErrorUrn.String()] = err

			//TODO: ONLY log the error if it is not a cancellation error.
			//This is the only condition so far when NOT to log the error!
			//logrus.Errorf("HTTP ERROR %v", err)
			return result
		}
		defer func() {
			err := response.Body.Close()
			if err != nil {
				logger.Log.Warnf("Error closing response Body %v", err)
			}
		}()

		requestBytes, _ := httputil.DumpRequest(req, true)
		responseBytes, _ := httputil.DumpResponse(response, true)

		if response.StatusCode >= 500 {
			result[core.ErrorUrn.String()] = fmt.Sprintf("Server Error %d", response.StatusCode)
		}

		result[RequestURLUrn.String()] = req.URL.String()
		result[core.BytesSentCountUrn.String()] = len(requestBytes)
		result[core.BytesReceivedCountUrn.String()] = len(responseBytes)

		for k, v := range response.Header {
			var key = RequestHeadersUrn.Name(k).String()
			result[key] = strings.Join(v, ",")
		}

		for k, v := range response.Header {
			var key = ResponseHeadersUrn.Name(k).String()
			result[key] = strings.Join(v, ",")
		}

		result[ResponseStatusUrn.String()] = response.StatusCode

		result[core.BytesSentUrn.String()] = string(requestBytes)
		result[core.BytesReceivedUrn.String()] = string(responseBytes)

		return result
	}
}
