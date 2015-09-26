package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
)


type RecordedRequest struct {
	request *http.Request
	body    string
}
type HttpRequestPredicate func(request RecordedRequest) bool

type RequestRecordingServer struct {
	requests []RecordedRequest
	port     int
	server   *httptest.Server
}

func CreateRequestRecordingServer(port int) *RequestRecordingServer {
	return &RequestRecordingServer{
		requests: []RecordedRequest{},
		port:     port,
	}
}

func (instance *RequestRecordingServer) Start() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body,err := ioutil.ReadAll(r.Body)
		if err != nil{
			panic(err)
		}
		instance.requests = append(instance.requests, RecordedRequest{
			request : r,
			body : string(body),
		})
	})
	instance.server = httptest.NewUnstartedServer(handler)
	instance.server.Listener, _ = net.Listen("tcp", ":"+strconv.Itoa(instance.port))
	instance.server.Start()
}

func (instance *RequestRecordingServer) Stop() {
	instance.server.Close()
}

func (instance *RequestRecordingServer) Clear() {
	instance.requests = []RecordedRequest{}
}

func (instance *RequestRecordingServer) Find(predicates ...HttpRequestPredicate) bool {

	for _, request := range instance.requests {
		results := make([]bool, len(predicates))
		for index, predicate := range predicates {
			results[index] = predicate(request)
		}
		thing := true
		for _, result := range results {
			if !result {
				thing = false
				break
			}
		}
		if thing {
			return thing
		}
	}
	return false

}

func RequestWithPath(path string) HttpRequestPredicate {
	return HttpRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.URL.Path == path
		if !result {
			Log.Println(fmt.Sprintf("path does not equal %s it equals %s", path, r.request.URL.Path))
		}
		return result
	})
}

func RequestWithMethod(method string) HttpRequestPredicate {
	return HttpRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.Method == method
		if !result {
			Log.Println("request method does not equal " + method)
		}
		return result
	})
}

func RequestWithHeader(key string, value string) HttpRequestPredicate {
	return HttpRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.Header.Get(key) == value
		if !result {
			Log.Println(fmt.Sprintf("request method does not contain header with key %s and value %s actual %s", key, value, r.request.Header.Get(key)))
		}
		return result
	})
}

func RequestWithBody(value string) HttpRequestPredicate {
	return HttpRequestPredicate(func(r RecordedRequest) bool {
		result := string(r.body) == value
		if !result {
			Log.Println(fmt.Sprintf("request body does not equal %s it equals %s", value, r.body))
		}
		return result
	})
}

func RequestWithQuerystring(value string) HttpRequestPredicate {
	return HttpRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.URL.RawQuery == value
		if !result {
			Log.Println("request query does not equal " + value + " | it equals " + r.request.URL.RawQuery)
		}
		return result
	})
}
