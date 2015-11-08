package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

//RequestFunc ...
type RequestFunc func() (*http.Request, error)

//RequestReader ...
type RequestReader struct {
	Requests []RequestFunc
}

//NewRequestReader ...
func NewRequestReader(filePath string) *RequestReader {
	file, err := os.Open(filePath)
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closes the file")
		}
	}()
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

//Size ...
func (instance *RequestReader) Size() int {
	return len(instance.Requests)
}

//Read ...
func (instance *RequestReader) Read(index int) RequestFunc {
	return instance.Requests[index]
}
