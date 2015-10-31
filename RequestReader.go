package main

import (
	"bufio"
	"net/http"
	"os"
)

type RequestFunc func() (*http.Request, error)

type RequestReader struct {
	Requests []RequestFunc
}

func NewRequestReader(filePath string) *RequestReader {
	file, err := os.Open(filePath)
	defer file.Close()
	check(err)
	requests := []RequestFunc{}
	requestAdapter := NewRequestAdapter()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		requestFunc := requestAdapter.Create(line)
		requests = append(requests, requestFunc)
		check(err)
	}

	return &RequestReader{
		Requests: requests,
	}
}

func (instance *RequestReader) Size() int {
	return len(instance.Requests)
}

func (instance *RequestReader) Read(index int) RequestFunc {
	return instance.Requests[index]
}
