package main

import (
	"bufio"
	"net/http"
	"os"
)

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

func (instance *RequestReader) Size() int {
	return len(instance.Requests)
}

func (instance *RequestReader) Read(index int) *http.Request {
	return instance.Requests[index]
}
