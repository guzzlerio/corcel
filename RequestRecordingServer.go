package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func (instance *RequestRecordingServer) Contains(predicates ...HttpRequestPredicate) bool {

	for _, request := range instance.requests {
		requestResult := false
		for _, predicate := range predicates {
			requestResult = predicate(request)
			if !requestResult {
				break
			}
		}
		if requestResult {
			return true
		}
	}
	return false

}

func RequestWithPath(path string) HttpRequestPredicate {
	return HttpRequestPredicate(func(request *http.Request) bool {
		return request.URL.Path == path
	})
}
