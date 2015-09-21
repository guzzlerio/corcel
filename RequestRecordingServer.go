package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"log"
)

type HttpRequestPredicate func(request *http.Request) bool

type RequestRecordingServer struct {
	requests []*http.Request
	port     int
	server   *httptest.Server
}

func CreateRequestRecordingServer(port int) *RequestRecordingServer {
	return &RequestRecordingServer{
		requests: []*http.Request{},
		port:     port,
	}
}

func (instance *RequestRecordingServer) Start() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instance.requests = append(instance.requests, r)
	})
	instance.server = httptest.NewUnstartedServer(handler)
	instance.server.Listener, _ = net.Listen("tcp", ":"+strconv.Itoa(instance.port))
	instance.server.Start()
}

func (instance *RequestRecordingServer) Stop() {
	instance.server.Close()
}

func (instance *RequestRecordingServer) Clear() {
	instance.requests = []*http.Request{}
}

func (instance *RequestRecordingServer) Contains(predicates ...HttpRequestPredicate) bool {

	for _, request := range instance.requests {
		results := make([]bool, len(predicates))
		for index, predicate := range predicates {
			results[index] = predicate(request)
		}
		thing := true
		for _, result := range results {
			if (!result) {
				thing = false
				break
			}
		}
		if (thing) {
			return thing
		}
	}
	return false

}

func RequestWithPath(path string) HttpRequestPredicate {
	return HttpRequestPredicate(func(request *http.Request) bool {
		result := request.URL.Path == path
		if !result {
			log.Println(fmt.Sprintf("path does not equal %s it equals %s", path, request.URL.Path))
		}
		return result
	})
}

func RequestWithMethod(method string) HttpRequestPredicate{
	return HttpRequestPredicate(func(request *http.Request) bool {
		result := request.Method == method
		if !result {
			log.Println("request method does not equal " + method)
		}
		return result
	})
}
