package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type RequestStream interface {
	Read() *http.Request
}

type SequentialStream struct {
	reader *RequestReader
}

func (instance *SequentialStream) Read() <-chan *http.Request {
	ch := make(chan *http.Request)
	go func() {
		for i := 0; i < len(instance.reader.Requests); i++ {
			request := instance.reader.Read(i)
			ch <- request
		}
		close(ch)
	}()
	return ch
}

func NewSequentialStream(reader *RequestReader) *SequentialStream {
	return &SequentialStream{
		reader: reader,
	}
}

type RequestReader struct {
	Requests []*http.Request
}

func NewRequestReader(filePath string) *RequestReader {
	file, err := os.Open(filePath)
	defer file.Close()
	check(err)
	requests := []*http.Request{}
	requestAdapter := NewRequestAdapter()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		request, err := requestAdapter.Create(line)
		requests = append(requests, request)
		check(err)
	}

	return &RequestReader{
		Requests: requests,
	}
}

func (instance *RequestReader) Read(index int) *http.Request {
	return instance.Requests[index]
}

func (instance *RequestReader) NewSequentialStream() *SequentialStream {
	return NewSequentialStream(instance)
}

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
