package main

import (
	"fmt"
	"net/http"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("RequestReader", func() {

	var list []string
	var reader *RequestReader

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

    Describe("RandomStream", func(){
        It("Reads randomly", func(){
		    requestSet1 := []*http.Request{}
            requestSet2 := []*http.Request{}

            stream1 := reader.NewRandomStream()
            for request := range stream1.Read() {
                requestSet1 = append(requestSet1, request)
            }

            stream2 := reader.NewRandomStream()
            for request := range stream2.Read() {
                requestSet2 = append(requestSet2, request)
            }
            Expect(ConcatRequestPaths(requestSet1)).ToNot(Equal(ConcatRequestPaths(requestSet2)))
        })
    })

	It("Single reader iterates over lines in a file", func() {
		requests := []*http.Request{}
		stream := reader.NewSequentialStream()
		for request := range stream.Read() {
			requests = append(requests, request)
		}
		Expect(len(requests)).To(Equal(len(list)))
	})

    for _, numberOfWorkers := range NUMBER_OF_WORKERS_TO_TEST {
		It(fmt.Sprintf("Multiple readers totalling %v iterates over lines in a file", numberOfWorkers), func() {
			var wg sync.WaitGroup
            var mutex = &sync.Mutex{}
			wg.Add(numberOfWorkers)
			requests := []*http.Request{}
			for i := 0; i < numberOfWorkers; i++ {
				go func() {
					stream := reader.NewSequentialStream()
					for request := range stream.Read() {
                        mutex.Lock()
						requests = append(requests, request)
                        mutex.Unlock()
					}
                    wg.Done()
				}()
			}
            wg.Wait()
			Expect(len(requests)).To(Equal(len(list) * numberOfWorkers))
		})
	}
})
