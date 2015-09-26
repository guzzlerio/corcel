package main

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestRecordingServer", func() {

	var (
		server *RequestRecordingServer
	)

	BeforeEach(func() {
		server = CreateRequestRecordingServer(8080)
	})

	Describe("Find", func() {

		var sampleUrl string
		var data string

		Describe("Single Request", func() {
			var request *http.Request

			BeforeEach(func() {
				data = "a=1&b=2"
				sampleUrl = "http://localhost:80/Fubar?"+data
				request,_ = http.NewRequest("GET", sampleUrl, bytes.NewBuffer([]byte(data)))
				request.Header.Set("Content-type", "application/json")
				server.requests = append(server.requests, RecordedRequest{
					request: request,
					body: data,
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
				Expect(server.Find(RequestWithHeader("Content-type","application/json"))).To(Equal(true))
			})

			It("Body", func(){
				request,_ = http.NewRequest("POST", sampleUrl, bytes.NewBuffer([]byte(data)))
				server.Clear()
				server.requests = append(server.requests, RecordedRequest{
					request: request,
					body: data,
				})
				Expect(server.Find(RequestWithBody(data))).To(Equal(true))
			})

			It("Querystring", func(){
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
				request,_ := http.NewRequest("GET", sampleUrl, nil)
				request.Header.Set("Content-type", "application/json")
				server.requests = append(server.requests, RecordedRequest{
					request:request,
				})

				data = "a=1&b=2"
				postRequest,_ := http.NewRequest("POST", sampleUrl, bytes.NewBuffer([]byte(data)))
				server.requests = append(server.requests, RecordedRequest{
					request: postRequest,
					body: data,
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
				Expect(server.Find(RequestWithHeader("Content-type","application/json"))).To(Equal(true))
			})

			It("Body", func(){
				Expect(server.Find(RequestWithBody(data))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

	})

})
