package main

import (
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestRecordingServer", func() {

	var (
		server *RequestRecordingServer
	)

	BeforeEach(func() {
		server = &RequestRecordingServer{
			requests: []*http.Request{},
		}
	})

	Describe("Find", func() {

		Describe("Single Request", func() {
			BeforeEach(func() {
				sampleUrl, _ := url.Parse("http://localhost:80/Fubar")
				request := &http.Request{
					URL: sampleUrl,
					Method: "GET",
					Header: map[string][]string{},
				}
				request.Header.Set("Content-type", "application/json")
				server.requests = append(server.requests, request)
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

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

		Describe("Multiple Requests", func() {
			BeforeEach(func() {
				sampleUrl, _ := url.Parse("http://localhost:80/Fubar")
				request := &http.Request{
					URL: sampleUrl,
					Method: "GET",
					Header: map[string][]string{},
				}
				request.Header.Set("Content-type", "application/json")
				server.requests = append(server.requests, request)
				server.requests = append(server.requests, &http.Request{
					URL: sampleUrl,
					Method: "POST",
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

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

	})

})
