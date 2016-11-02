package request

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	. "ci.guzzler.io/guzzler/corcel/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestStream", func() {
	var (
		list     []string
		reader   *Reader
		iterator RequestStream
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

	Describe("Sequential Request Stream", func() {
		BeforeEach(func() {
			iterator = NewSequentialRequestStream(reader)
		})

		It("executes in a sequential order", func() {
			actual := []*http.Request{}
			for iterator.HasNext() {
				req, _ := iterator.Next()
				actual = append(actual, req)
			}
			Expect(len(actual)).To(Equal(len(list)))
		})

		It("calculates Size", func() {
			Expect(iterator.Size()).To(Equal(len(list)))
		})

		It("calculates Progress", func() {
			for i := 0; i < (len(list) / 2); i++ {
				_, _ = iterator.Next()
			}
			Expect(iterator.Progress()).To(Equal(50))
		})
	})

	Describe("Timebased RequestStream", func() {
		var duration time.Duration
		var lock = &sync.Mutex{}

		BeforeEach(func() {
			lock.Lock()
			iterator = NewSequentialRequestStream(reader)
			duration = time.Duration(1 * time.Second)
			iterator = NewTimeBasedRequestStream(iterator, duration)
			lock.Unlock()
		})

		It("executes for the duration", func() {
			actual := Time(func() {
				for iterator.HasNext() {
					_, err := iterator.Next()
					check(err)
				}
			})
			max := time.Duration(duration + (500 * time.Millisecond))
			Expect(DurationIsBetween(actual, duration, max)).To(Equal(true))
		})

		It("calculates Progress", func() {
			go func() {
				lock.Lock()
				for iterator.HasNext() {
					_, _ = iterator.Next()
				}
				lock.Unlock()
			}()
			time.Sleep((100 * time.Millisecond))
			Expect(iterator.Progress()).To(BeNumerically(">=", 10))
			Expect(iterator.Progress()).To(BeNumerically("<", 20))
		})

		It("calculates Size", func() {
			Expect(iterator.Size()).To(Equal(int(duration.Nanoseconds())))
		})
	})

	Describe("Random RequestStream", func() {
		It("calculates Size", func() {
			stream := NewRandomRequestStream(reader)
			Expect(stream.Size()).To(Equal(len(list)))
		})

		It("executes in a random order", func() {
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

		It("from Reader with size of 1", func() {
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

})
