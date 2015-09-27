package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestRecordingServer", func() {

	var (
		server *RequestRecordingServer
		port   int
	)

	BeforeEach(func() {
		port = 8080
		server = CreateRequestRecordingServer(port)
	})

	Describe("Find", func() {

		var sampleUrl string
		var data string

		Describe("Single Request", func() {
			var request *http.Request

			BeforeEach(func() {
				data = "a=1&b=2"
				sampleUrl = "http://localhost:80/Fubar?" + data
				request, _ = http.NewRequest("GET", sampleUrl, bytes.NewBuffer([]byte(data)))
				request.Header.Set("Content-type", "application/json")
				server.requests = append(server.requests, RecordedRequest{
					request: request,
					body:    data,
				})
			})

			It("Path", func() {
				expectedPath := "/Fubar"
				Expect(server.Find(RequestWithPath(expectedPath))).To(Equal(true))
			})

			It("Method", func() {
				expectedMethod := "GET"
				Expect(server.Find(RequestWithMethod(expectedMethod))).To(Equal(true))
			})

			It("Header", func() {
				Expect(server.Find(RequestWithHeader("Content-type", "application/json"))).To(Equal(true))
			})

			It("Body", func() {
				request, _ = http.NewRequest("POST", sampleUrl, bytes.NewBuffer([]byte(data)))
				server.Clear()
				server.requests = append(server.requests, RecordedRequest{
					request: request,
					body:    data,
				})
				Expect(server.Find(RequestWithBody(data))).To(Equal(true))
			})

			It("Querystring", func() {
				Expect(server.Find(RequestWithQuerystring(data))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

		Describe("Multiple Requests", func() {
			BeforeEach(func() {
				sampleUrl = "http://localhost:80/Fubar"
				request, _ := http.NewRequest("GET", sampleUrl, nil)
				request.Header.Set("Content-type", "application/json")
				server.requests = append(server.requests, RecordedRequest{
					request: request,
				})

				data = "a=1&b=2"
				postRequest, _ := http.NewRequest("POST", sampleUrl, bytes.NewBuffer([]byte(data)))
				server.requests = append(server.requests, RecordedRequest{
					request: postRequest,
					body:    data,
				})
			})

			It("Path", func() {
				expectedPath := "/Fubar"
				Expect(server.Find(RequestWithPath(expectedPath))).To(Equal(true))
			})

			It("Method", func() {
				expectedMethod := "GET"
				Expect(server.Find(RequestWithMethod(expectedMethod))).To(Equal(true))
			})

			It("Header", func() {
				Expect(server.Find(RequestWithHeader("Content-type", "application/json"))).To(Equal(true))
			})

			It("Body", func() {
				Expect(server.Find(RequestWithBody(data))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

	})

	Describe("Response factory", func() {
		BeforeEach(func(){
			server.Start()
		})

		AfterEach(func(){
			server.Clear()
			server.Stop()
		})

		It("Defines the response to be used for the server", func() {
			message := "Hello World"

			factory := HttpResponseFactory(func(w http.ResponseWriter) {
				io.WriteString(w, message)
			})

			server.Use(factory)

			response, body, err := HttpRequestDo("GET", fmt.Sprintf("http://localhost:%d", port), nil, nil)

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(string(body)).To(Equal(message))
		})

		It("Clears the response to be used for the server", func() {
			message := "Hello World"
			factory := HttpResponseFactory(func(w http.ResponseWriter) {
				io.WriteString(w, message)
			})
			server.Use(factory)
			server.Clear()

			response, body, err := HttpRequestDo("GET", fmt.Sprintf("http://localhost:%d", port), nil, nil)

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(string(body)).To(Equal(""))
		})

		It("Defines the response to be used for the server with predicate", func() {
			message := "Hello World"
			factory := HttpResponseFactory(func(w http.ResponseWriter) {
				io.WriteString(w, message)
			})

			predicates := []HttpRequestPredicate{
				RequestWithPath("/talula"),
				RequestWithMethod("POST"),
				RequestWithHeader("Content-Type", "application/json"),
			}

			server.Use(factory).For(predicates...)

			pathMatching := fmt.Sprintf("http://localhost:%d/talula", port)
			verbMatching := "POST"
			responseMatching, bodyMatching, errMatching := HttpRequestDo(verbMatching, pathMatching , nil, func(request *http.Request){
				request.Header.Set("Content-Type", "application/json")
			})

			pathNonMatching := fmt.Sprintf("http://localhost:%d", port)
			verbNonMatching := "GET"
			responseNonMatching, bodyNonMatching, errNonMatching := HttpRequestDo(verbNonMatching, pathNonMatching , nil, nil)

			Expect(errMatching).To(BeNil())
			Expect(responseMatching.StatusCode).To(Equal(http.StatusOK))
			Expect(string(bodyMatching)).To(Equal(message))

			Expect(errNonMatching).To(BeNil())
			Expect(responseNonMatching.StatusCode).To(Equal(http.StatusOK))
			Expect(string(bodyNonMatching)).To(Equal(""))
		})
	})
})
