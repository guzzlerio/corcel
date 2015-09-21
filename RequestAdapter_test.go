package main

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestAdapter", func() {
	var (
		url     string
		line    string
		adapter RequestAdapter
		req     *http.Request
		err     error
	)

	BeforeEach(func() {
		url = "http://localhost:8000/A"
		line = url
		line += " -X POST"
		adapter = NewRequestAdapter()
		req, err = adapter.Create(line)
		Expect(err).To(BeNil())
	})

	It("Parses URL", func() {
		Expect(req.URL.String()).To(Equal(url))
	})

	It("Parses Method", func() {
		Expect(req.Method).To(Equal("POST"))
	})

	It("Parses URLs with leading whitespace", func() {
		line = "      " + url
		adapter = NewRequestAdapter()
		req, err = adapter.Create(line)
		Expect(req.URL.String()).To(Equal(url))
	})
})
