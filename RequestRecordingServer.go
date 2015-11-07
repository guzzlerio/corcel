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

type RecordedRequest struct {
	request *http.Request
	body    string
}

type HttpRequestPredicate func(request RecordedRequest) bool
type HttpResponseFactory func(writer http.ResponseWriter)

type UseWithPredicates struct {
	ResponseFactory   HttpResponseFactory
	RequestPredicates []HttpRequestPredicate
}

type RequestRecordingServer struct {
	requests []RecordedRequest
	port     int
	server   *httptest.Server
	use      []UseWithPredicates
}

func CreateRequestRecordingServer(port int) *RequestRecordingServer {
	return &RequestRecordingServer{
		requests: []RecordedRequest{},
		port:     port,
		use:      []UseWithPredicates{},
	}
}

func (instance *RequestRecordingServer) Start() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		check(err)
		recordedRequest := RecordedRequest{
			request: r,
			body:    string(body),
		}
		instance.requests = append(instance.requests,recordedRequest)
		if instance.use != nil {
			for _, item := range instance.use {
				if item.RequestPredicates != nil {
					result := instance.Evaluate(recordedRequest, item.RequestPredicates...)
					if (result){
						item.ResponseFactory(w)
						return
					}
				}else{
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

func (instance *RequestRecordingServer) Stop() {
	instance.server.Close()
	time.Sleep(1 * time.Microsecond)
}

func (instance *RequestRecordingServer) Clear() {
	instance.requests = []RecordedRequest{}
	instance.use = []UseWithPredicates{}
}

func (instance *RequestRecordingServer) Evaluate(request RecordedRequest, predicates ...HttpRequestPredicate) bool{
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

func (instance *RequestRecordingServer) Find(predicates ...HttpRequestPredicate) bool {
	for _, request := range instance.requests {
		if instance.Evaluate(request) {
			return true
		}
	}
	return false
}

func (instance *RequestRecordingServer) Use(factory HttpResponseFactory) *RequestRecordingServer{
	instance.use = append(instance.use, UseWithPredicates{
		ResponseFactory : factory,
		RequestPredicates : []HttpRequestPredicate{},
	})
	return instance
}

func (instance *RequestRecordingServer) For(predicates ...HttpRequestPredicate){
	index := len(instance.use) - 1
	for _, item := range predicates {
		instance.use[index].RequestPredicates = append(instance.use[index].RequestPredicates, item)
	}
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
