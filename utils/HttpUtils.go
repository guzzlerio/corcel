package utils

import (
	"io"
	"io/ioutil"
	"net/http"

	"ci.guzzler.io/guzzler/corcel/errormanager"
)

//HTTPRequestDo ...
func HTTPRequestDo(verb string, url string, bodyBuffer io.Reader, changeRequestDelegate func(request *http.Request)) (response *http.Response, body string, err error) {
	client := &http.Client{}
	request, err := http.NewRequest(verb, url, bodyBuffer)
	errormanager.Check(err)
	if changeRequestDelegate != nil {
		changeRequestDelegate(request)
	}
	response, err = client.Do(request)
	errormanager.Check(err)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err == nil {
		body = string(bodyBytes)
	}
	return
}

//ConcatRequestPaths ...
func ConcatRequestPaths(requests []*http.Request) string {
	result := ""
	for _, request := range requests {
		result = result + request.URL.Path
	}
	return result
}