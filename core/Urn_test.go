package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Urn", func() {
	It("can create a counter", func() {
		urn := NewUrn("http").Counter().Name("status", "all").Name("200").String()
		Expect(urn).To(Equal("urn:http:counter:status:all:200"))
	})

	It("can create a gauge", func() {
		urn := NewUrn("http").Gauge().Name("status", "all").Name("200").String()
		Expect(urn).To(Equal("urn:http:gauge:status:all:200"))
	})

	It("can create a meter", func() {
		urn := NewUrn("http").Meter().Name("status", "all").Name("200").String()
		Expect(urn).To(Equal("urn:http:meter:status:all:200"))
	})

	It("can create a timer", func() {
		urn := NewUrn("http").Timer().Name("status", "all").Name("200").String()
		Expect(urn).To(Equal("urn:http:timer:status:all:200"))
	})

	It("can create a histogram", func() {
		urn := NewUrn("http").Timer().Name("status", "all").Name("200").String()
		Expect(urn).To(Equal("urn:http:timer:status:all:200"))
	})

	It("ignores the metric is not specified", func() {
		urn := NewUrn("http").Name("status").Name("200").String()
		Expect(urn).To(Equal("urn:http:status:200"))
	})

	It("builds the urn in lowercase", func() {

	})
})
