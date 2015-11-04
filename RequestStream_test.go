package main

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestStream", func() {
	var (
		list   []string
		reader *RequestReader
	)

	BeforeEach(func() {
		list = []string{
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/1")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/2")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/3")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/4")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/5")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/6")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/7")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/8")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/9")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/10")),
		}
		file := CreateFileFromLines(list)
		file.Close()
		reader = NewRequestReader(file.Name())
	})

	It("Sequential Request Stream", func() {
		iterator := NewSequentialRequestStream(reader)

		actual := []*http.Request{}
		for iterator.HasNext() {
            req,_ := iterator.Next()
			actual = append(actual, req)
		}
		Expect(len(actual)).To(Equal(len(list)))
	})

	It("Random Request Stream", func() {
		requestSet1 := []*http.Request{}
		requestSet2 := []*http.Request{}

		stream1 := NewRandomRequestStream(reader)
		for stream1.HasNext() {
            req,_ := stream1.Next()
			requestSet1 = append(requestSet1, req)
		}

		stream2 := NewRandomRequestStream(reader)
		for stream2.HasNext() {
            req,_ := stream2.Next()
			requestSet2 = append(requestSet2, req)
		}

		Expect(len(requestSet1)).To(Equal(len(list)))
		Expect(len(requestSet2)).To(Equal(len(list)))
		Expect(ConcatRequestPaths(requestSet1)).ToNot(Equal(ConcatRequestPaths(requestSet2)))
	})

	It("Timebased RequestStream", func() {
		iterator := NewSequentialRequestStream(reader)
		duration := time.Duration(3 * time.Second)
		iterator = NewTimeBasedRequestStream(iterator, duration)
		actual := Time(func() {
			for iterator.HasNext() {
				iterator.Next()
			}
		})
		max := time.Duration(duration + (500 * time.Millisecond))
		Expect(DurationIsBetween(actual, duration, max)).To(Equal(true))
	})

    It("Random Request Stream from Reader with size of 1", func(){
		list = []string{
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/1")),
		}
		file := CreateFileFromLines(list)
		file.Close()
		reader = NewRequestReader(file.Name())

		requestSet := []*http.Request{}

		stream := NewRandomRequestStream(reader)
		for stream.HasNext() {
            req,_ := stream.Next()
			requestSet = append(requestSet, req)
		}
		Expect(len(requestSet)).To(Equal(len(list)))
    })

})
