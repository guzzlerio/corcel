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
	client  *http.Client
	URL     string
	Method  string
	Body    string
	Headers http.Header
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
		client: &http.Client{Transport: tr},
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

		if instance.Body != "" && instance.Body[0] == '@' {
			contents, err := ioutil.ReadFile(instance.Body[1:])
			if err != nil {
				result[core.ErrorUrn.String()] = err
				return result
			}
			instance.Body = string(contents)
		}

		var requestURL = instance.URL
		var method = instance.Method
		var headers = http.Header{}
		var body = instance.Body

		for k := range instance.Headers {
			headers.Set(k, instance.Headers.Get(k))
		}
		if executionContext["$httpHeaders"] != nil {
			for hKey, hValue := range executionContext["$httpHeaders"].(map[string]interface{}) {
				headerKey := hKey
				headerValue := hValue.(string)

				if headers.Get(headerKey) == "" {
					headers.Set(headerKey, headerValue)
				}
			}
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
		//result[RequestHeadersUrn.String()] = req.Header

		for k, v := range response.Header {
			var key = RequestHeadersUrn.Name(k).String()
			result[key] = strings.Join(v, ",")
		}

		for k, v := range response.Header {
			var key = ResponseHeadersUrn.Name(k).String()
			result[key] = strings.Join(v, ",")
		}

		//TODO: We need a Response Headers key too
		result[ResponseStatusUrn.String()] = response.StatusCode

		result[core.BytesSentUrn.String()] = string(requestBytes)
		result[core.BytesReceivedUrn.String()] = string(responseBytes)

		return result
	}
}
