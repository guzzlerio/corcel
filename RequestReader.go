package main

import (
	"bufio"
	"net/http"
	"os"

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
