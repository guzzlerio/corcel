package utils

import (
	"net/http"

	"github.com/guzzlerio/rizo"
)

//ConcatRequestPaths ...
func ConcatRequestPaths(requests []*http.Request) string {
	result := ""
	for _, request := range requests {
		result = result + request.URL.Path
	}
	return result
}

//ToHTTPRequestArray ...
func ToHTTPRequestArray(inArray []rizo.RecordedRequest) []*http.Request {
	returnArray := []*http.Request{}
	for _, req := range inArray {
		returnArray = append(returnArray, req.Request)
	}
	return returnArray
}
