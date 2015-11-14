package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

//HTTPRequestDo ...
func HTTPRequestDo(verb string, url string, bodyBuffer io.Reader, changeRequestDelegate func(request *http.Request)) (response *http.Response, body string, err error) {
	client := &http.Client{}
	request, err := http.NewRequest(verb, url, bodyBuffer)
	check(err)
	if changeRequestDelegate != nil {
		changeRequestDelegate(request)
	}
	response, err = client.Do(request)
	check(err)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err == nil {
		body = string(bodyBytes)
	}
	return
}
