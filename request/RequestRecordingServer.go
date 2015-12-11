package request

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"ci.guzzler.io/guzzler/corcel/logger"
)

//RecordedRequest ...
type RecordedRequest struct {
	Request *http.Request
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
	Requests []RecordedRequest
	port     int
	server   *httptest.Server
	use      []UseWithPredicates
}

//CreateRequestRecordingServer ...
func CreateRequestRecordingServer(port int) *RequestRecordingServer {
	return &RequestRecordingServer{
		Requests: []RecordedRequest{},
		port:     port,
		use:      []UseWithPredicates{},
	}
}

func (instance *RequestRecordingServer) evaluatePredicates(recordedRequest RecordedRequest, w http.ResponseWriter) {
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
}

//Start ...
func (instance *RequestRecordingServer) Start() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		check(err)
		recordedRequest := RecordedRequest{
			Request: r,
			body:    string(body),
		}
		instance.Requests = append(instance.Requests, recordedRequest)
		if instance.use != nil {
			instance.evaluatePredicates(recordedRequest, w)
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
	instance.Requests = []RecordedRequest{}
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
	for _, request := range instance.Requests {
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
		result := r.Request.URL.Path == path
		if !result {
			logger.Log.Println(fmt.Sprintf("path does not equal %s it equals %s", path, r.Request.URL.Path))
		}
		return result
	})
}

//RequestWithMethod ...
func RequestWithMethod(method string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.Method == method
		if !result {
			logger.Log.Println("request method does not equal " + method)
		}
		return result
	})
}

//RequestWithHeader ...
func RequestWithHeader(key string, value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.Header.Get(key) == value
		if !result {
			logger.Log.Println(fmt.Sprintf("request method does not contain header with key %s and value %s actual %s", key, value, r.Request.Header.Get(key)))
		}
		return result
	})
}

//RequestWithBody ...
func RequestWithBody(value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := string(r.body) == value
		if !result {
			logger.Log.Println(fmt.Sprintf("request body does not equal %s it equals %s", value, r.body))
		}
		return result
	})
}

//RequestWithQuerystring ...
func RequestWithQuerystring(value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.URL.RawQuery == value
		if !result {
			logger.Log.Println("request query does not equal " + value + " | it equals " + r.Request.URL.RawQuery)
		}
		return result
	})
}
