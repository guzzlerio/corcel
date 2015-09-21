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

	Describe("Contains", func() {

		Describe("Single Request", func() {
			BeforeEach(func() {
				sampleUrl, _ := url.Parse("http://localhost:80/Fubar")
				server.requests = append(server.requests, &http.Request{
					URL: sampleUrl,
					Method: "GET",
				})
			})

			It("Path", func() {
				expectedPath := "/Fubar"
				Expect(server.Contains(RequestWithPath(expectedPath))).To(Equal(true))
			})

			It("Method", func() {
				expectedMethod := "GET"
				Expect(server.Contains(RequestWithMethod(expectedMethod))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Contains(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

		Describe("Multiple Requests", func() {
			BeforeEach(func() {
				sampleUrl, _ := url.Parse("http://localhost:80/Fubar")
				server.requests = append(server.requests, &http.Request{
					URL: sampleUrl,
					Method: "GET",
				})
				server.requests = append(server.requests, &http.Request{
					URL: sampleUrl,
					Method: "POST",
				})
			})

			It("Path", func() {
				expectedPath := "/Fubar"
				Expect(server.Contains(RequestWithPath(expectedPath))).To(Equal(true))
			})

			It("Method", func() {
				expectedMethod := "GET"
				Expect(server.Contains(RequestWithMethod(expectedMethod))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(server.Contains(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

	})

})
