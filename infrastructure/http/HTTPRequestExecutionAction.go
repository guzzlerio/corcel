package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/logger"
)

//HTTPRequestExecutionAction ...
type HTTPRequestExecutionAction struct {
	Client  *http.Client
	URL     string
	Method  string
	Body    string
	Headers http.Header
}

func (instance *HTTPRequestExecutionAction) initialize() {
	instance.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}
}

//Execute ...
func (instance *HTTPRequestExecutionAction) Execute(context core.ExecutionContext, cancellation chan struct{}) core.ExecutionResult {
	if instance.Client == nil {
		instance.initialize()
	}

	result := core.ExecutionResult{}

	if instance.Body != "" && instance.Body[0] == '@' {
		contents, err := ioutil.ReadFile(instance.Body[1:])
		if err != nil {
			result["action:error"] = err
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
	if context["httpHeaders"] != nil {
		for hKey, hValue := range context["httpHeaders"].(map[interface{}]interface{}) {
			headerKey := hKey.(string)
			headerValue := hValue.(string)

			if headers.Get(headerKey) == "" {
				headers.Set(headerKey, headerValue)
			}
		}
	}

	//TODO:  This is inefficient but working.  We can always make this
	// better and speed it up later!!
	//
	//TODO: Replace with some efficient templating engine.  Mustache
	// might even work who knows
	for k, v := range context {
		token := "$" + k
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
	req.Cancel = cancellation
	//This should be a configuration item.  It allows the client to work
	//in a way similar to a server which does not support HTTP KeepAlive
	//After each request the client channel is closed.  When set to true
	//the performance overhead is large in terms of Network IO throughput

	//req.Close = true

	if err != nil {
		result["action:error"] = err
		return result
	}

	req.Header = headers

	response, err := instance.Client.Do(req)
	if err != nil {
		result["action:error"] = err
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
		result["action:error"] = fmt.Sprintf("Server Error %d", response.StatusCode)
	}

	result["http:request:url"] = req.URL.String()
	result["action:bytes:sent"] = len(requestBytes)
	result["action:bytes:received"] = len(responseBytes)
	result["http:request:headers"] = req.Header
	result["http:response:status"] = response.StatusCode
	result["http:response:body"] = string(responseBytes)

	return result
}
