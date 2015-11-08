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
	if err != nil {
		Log.Printf("err creating request %v")
		return
	}
	if changeRequestDelegate != nil {
		changeRequestDelegate(request)
	}
	response, err = client.Do(request)
	if err != nil {
		Log.Printf("err getting response %v")
		return
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log.Printf("err reading body %v")
	} else {
		body = string(bodyBytes)
	}
	return
}
