package http

import (
	nethttp "net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Http DefaultsExtractor", func() {
	Describe("Sets Defaults", func() {
		var result HttpActionState
		BeforeEach(func() {
			var input = map[string]interface{}{}
			var extractor = NewDefaultsExtractor()

			result = extractor.Extract(input)
		})
		It("sets default URL to empty string", func() {
			Expect(result.URL).To(Equal(""))
		})
		It("sets default method to GET", func() {
			Expect(result.Method).To(Equal("GET"))
		})
		It("sets default body to empty string", func() {
			Expect(result.Body).To(Equal(""))
		})
		It("sets default header to empty collection", func() {
			Expect(result.Headers).To(Equal(nethttp.Header{}))
		})
	})

	It("Extracts HttpActionState", func() {
		var input = map[string]interface{}{}
		input["defaults"] = map[string]interface{}{}
		var defaults = input["defaults"].(map[string]interface{})
		defaults["HttpAction"] = map[string]interface{}{}

		var action = defaults["HttpAction"].(map[string]interface{})
		action["headers"] = map[string]interface{}{}

		var headers = action["headers"].(map[string]interface{})

		headers["key"] = "value"

		action["method"] = "GET"
		action["body"] = "Bang Bang"
		action["url"] = "http://somewhere"

		var extractor = NewDefaultsExtractor()

		var result = extractor.Extract(input)

		Expect(result.Headers.Get("key")).To(Equal("value"))
		Expect(result.Method).To(Equal(action["method"]))
		Expect(result.Body).To(Equal(action["body"]))
		Expect(result.URL).To(Equal(action["url"]))
	})

	It("returns Empty state when no defaults", func() {
		var extractor = NewDefaultsExtractor()
		var state = map[string]interface{}{}
		Expect(extractor.Extract(state)).ToNot(BeNil())
	})

	It("returns Empty state when no default HttpAction", func() {
		var extractor = NewDefaultsExtractor()
		var state = map[string]interface{}{}
		state["defaults"] = map[string]interface{}{}
		Expect(extractor.Extract(state)).ToNot(BeNil())
	})

	It("returns Empty state when no http definitions", func() {
		var extractor = NewDefaultsExtractor()
		var state = map[string]interface{}{}
		state["defaults"] = map[string]interface{}{}
		var defaults = state["defaults"].(map[string]interface{})
		defaults["HttpAction"] = map[string]interface{}{}
		Expect(extractor.Extract(state)).ToNot(BeNil())
	})
})
