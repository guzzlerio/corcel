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
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
		}
		file := CreateFileFromLines(list)
		file.Close()
		reader = NewRequestReader(file.Name())
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
