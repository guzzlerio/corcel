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
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/1")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/2")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/3")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/4")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/5")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/6")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/7")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/8")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/9")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/10")),
		}
		file := CreateFileFromLines(list)
		check(file.Close())
		reader = NewRequestReader(file.Name())
	})

	It("Sequential Request Stream", func() {
		iterator := NewSequentialRequestStream(reader)

		actual := []*http.Request{}
		for iterator.HasNext() {
			req, _ := iterator.Next()
			actual = append(actual, req)
		}
		Expect(len(actual)).To(Equal(len(list)))
	})

	It("Random Request Stream", func() {
		requestSet1 := []*http.Request{}
		requestSet2 := []*http.Request{}

		stream1 := NewRandomRequestStream(reader)
		for stream1.HasNext() {
			req, _ := stream1.Next()
			requestSet1 = append(requestSet1, req)
		}

		stream2 := NewRandomRequestStream(reader)
		for stream2.HasNext() {
			req, _ := stream2.Next()
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
				_, err := iterator.Next()
				check(err)
			}
		})
		max := time.Duration(duration + (500 * time.Millisecond))
		Expect(DurationIsBetween(actual, duration, max)).To(Equal(true))
	})

	It("Random Request Stream from Reader with size of 1", func() {
		list = []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/1")),
		}
		file := CreateFileFromLines(list)
		check(file.Close())
		reader = NewRequestReader(file.Name())

		requestSet := []*http.Request{}

		stream := NewRandomRequestStream(reader)
		for stream.HasNext() {
			req, _ := stream.Next()
			requestSet = append(requestSet, req)
		}
		Expect(len(requestSet)).To(Equal(len(list)))
	})

})
