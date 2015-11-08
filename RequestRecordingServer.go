package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

//RecordedRequest ...
type RecordedRequest struct {
	request *http.Request
	body    string
}

//HTTPRequestPredicate ...
type HTTPRequestPredicate func(request RecordedRequest) bool

//HTTPResponseFactory ...
type HTTPResponseFactory func(writer http.ResponseWriter)

//UseWithPredicates ...
type UseWithPredicates struct {
	ResponseFactory   HTTPResponseFactory
	RequestPredicates []HTTPRequestPredicate
}

//RequestRecordingServer ...
type RequestRecordingServer struct {
	requests []RecordedRequest
	port     int
	server   *httptest.Server
	use      []UseWithPredicates
}

//CreateRequestRecordingServer ...
func CreateRequestRecordingServer(port int) *RequestRecordingServer {
	return &RequestRecordingServer{
		requests: []RecordedRequest{},
		port:     port,
		use:      []UseWithPredicates{},
	}
}

//Start ...
func (instance *RequestRecordingServer) Start() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		check(err)
		recordedRequest := RecordedRequest{
			request: r,
			body:    string(body),
		}
		instance.requests = append(instance.requests, recordedRequest)
		if instance.use != nil {
			for _, item := range instance.use {
				if item.RequestPredicates != nil {
					result := instance.Evaluate(recordedRequest, item.RequestPredicates...)
					if result {
						item.ResponseFactory(w)
						return
					}
				} else {
					item.ResponseFactory(w)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
	instance.server = httptest.NewUnstartedServer(handler)
	instance.server.Listener, _ = net.Listen("tcp", ":"+strconv.Itoa(instance.port))
	instance.server.Start()
}

//Stop ...
func (instance *RequestRecordingServer) Stop() {
	instance.server.Close()
	time.Sleep(1 * time.Microsecond)
}

//Clear ...
func (instance *RequestRecordingServer) Clear() {
	instance.requests = []RecordedRequest{}
	instance.use = []UseWithPredicates{}
}

//Evaluate ...
func (instance *RequestRecordingServer) Evaluate(request RecordedRequest, predicates ...HTTPRequestPredicate) bool {
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
	return thing
}

//Find ...
func (instance *RequestRecordingServer) Find(predicates ...HTTPRequestPredicate) bool {
	for _, request := range instance.requests {
		if instance.Evaluate(request) {
			return true
		}
	}
	return false
}

//Use ...
func (instance *RequestRecordingServer) Use(factory HTTPResponseFactory) *RequestRecordingServer {
	instance.use = append(instance.use, UseWithPredicates{
		ResponseFactory:   factory,
		RequestPredicates: []HTTPRequestPredicate{},
	})
	return instance
}

//For ...
func (instance *RequestRecordingServer) For(predicates ...HTTPRequestPredicate) {
	index := len(instance.use) - 1
	for _, item := range predicates {
		instance.use[index].RequestPredicates = append(instance.use[index].RequestPredicates, item)
	}
}

//RequestWithPath ...
func RequestWithPath(path string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.URL.Path == path
		if !result {
			Log.Println(fmt.Sprintf("path does not equal %s it equals %s", path, r.request.URL.Path))
		}
		return result
	})
}

//RequestWithMethod ...
func RequestWithMethod(method string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.Method == method
		if !result {
			Log.Println("request method does not equal " + method)
		}
		return result
	})
}

//RequestWithHeader ...
func RequestWithHeader(key string, value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.Header.Get(key) == value
		if !result {
			Log.Println(fmt.Sprintf("request method does not contain header with key %s and value %s actual %s", key, value, r.request.Header.Get(key)))
		}
		return result
	})
}

//RequestWithBody ...
func RequestWithBody(value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := string(r.body) == value
		if !result {
			Log.Println(fmt.Sprintf("request body does not equal %s it equals %s", value, r.body))
		}
		return result
	})
}

//RequestWithQuerystring ...
func RequestWithQuerystring(value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.request.URL.RawQuery == value
		if !result {
			Log.Println("request query does not equal " + value + " | it equals " + r.request.URL.RawQuery)
		}
		return result
	})
}
