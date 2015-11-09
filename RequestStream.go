package main

import (
	"net/http"
	"time"
)

//RequestStream ...
type RequestStream interface {
	HasNext() bool
	Next() (*http.Request, error)
	Reset()
}

//SequentialRequestStream ...
type SequentialRequestStream struct {
	reader  *RequestReader
	current int
}

//NewSequentialRequestStream ...
func NewSequentialRequestStream(reader *RequestReader) RequestStream {
	return &SequentialRequestStream{
		reader:  reader,
		current: 0,
	}
}

//HasNext ...
func (instance *SequentialRequestStream) HasNext() bool {
	return instance.current < instance.reader.Size()
}

//Next ...
func (instance *SequentialRequestStream) Next() (*http.Request, error) {
	element := instance.reader.Read(instance.current)
	instance.current++
	return element()
}

//Reset ...
func (instance *SequentialRequestStream) Reset() {
	instance.current = 0
}

//RandomRequestStream ...
type RandomRequestStream struct {
	reader *RequestReader
	count  int
}

//NewRandomRequestStream ...
func NewRandomRequestStream(reader *RequestReader) RequestStream {
	return &RandomRequestStream{
		reader: reader,
		count:  0,
	}
}

//HasNext ...
func (instance *RandomRequestStream) HasNext() bool {
	return instance.count < instance.reader.Size()
}

//Next ...
func (instance *RandomRequestStream) Next() (*http.Request, error) {
	if instance.reader.Size() == 0 {
		panic("The reader is empty")
	}
	max := instance.reader.Size() - 1
	if max == 0 {
		max = 1
	}
	randomIndex := Random.Intn(max)
	element := instance.reader.Read(randomIndex)
	instance.count++
	return element()
}

//Reset ...
func (instance *RandomRequestStream) Reset() {
	instance.count = 0
}

//TimeBasedRequestStream ...
type TimeBasedRequestStream struct {
	stream   RequestStream
	duration time.Duration
	start    time.Time
}

//NewTimeBasedRequestStream ...
func NewTimeBasedRequestStream(stream RequestStream, duration time.Duration) RequestStream {
	return &TimeBasedRequestStream{
		stream:   stream,
		duration: duration,
	}
}

//HasNext ...
func (instance *TimeBasedRequestStream) HasNext() bool {
	if instance.start.IsZero() {
		instance.start = time.Now()
	}
	return time.Since(instance.start) < instance.duration
}

//Next ...
func (instance *TimeBasedRequestStream) Next() (*http.Request, error) {
	if !instance.stream.HasNext() {
		instance.stream.Reset()
	}
	return instance.stream.Next()
}

//Reset ...
func (instance *TimeBasedRequestStream) Reset() {
	instance.start = time.Time{}
}
