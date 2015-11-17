package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//RequestAdapter ...
type RequestAdapter struct{}

//NewRequestAdapter ...
func NewRequestAdapter() RequestAdapter {
	return RequestAdapter{}
}

type RequestConfigHandler interface {
	Handle(options []string, index int, req *http.Request) *http.Request
}

func HandlerForMethod(options []string, index int, req *http.Request) *http.Request {
	req.Method = options[index+1]
	return req
}

//Create ...
func (instance RequestAdapter) Create(line string) RequestFunc {
	return RequestFunc(func() (*http.Request, error) {
		commandLineLexer := NewCommandLineLexer()
		lineSplit := commandLineLexer.Lex(line)
		req, err := http.NewRequest("GET", lineSplit[0], nil)
		if err != nil {
			return nil, err
		}
		for index := range lineSplit {
			if lineSplit[index] == "-X" {
				//req.Method = lineSplit[index+1]
				req = HandlerForMethod(lineSplit, index, req)
			}
			if lineSplit[index] == "-H" {
				value := strings.Trim(lineSplit[index+1], "\"")

				valueSplit := strings.Split(value, ":")
				req.Header.Set(strings.TrimSpace(valueSplit[0]), strings.TrimSpace(valueSplit[1]))
			}
			if lineSplit[index] == "-d" {
				rawBody := lineSplit[index+1]

				if strings.ToLower(req.Method) == "get" {
					req.URL.RawQuery = lineSplit[index+1]
				} else {
					var body *bytes.Buffer
					bodyBytes := []byte(rawBody)
					if strings.HasPrefix(rawBody, "@") {
						body = loadRequestBodyFromFile(string(bytes.TrimLeft(bodyBytes, "@")))
					} else {
						Log.Println("body from request")
						body = bytes.NewBuffer(bodyBytes)
					}
					req, err = http.NewRequest(req.Method, req.URL.String(), body)
				}
			}
			if lineSplit[index] == "-A" {
				req.Header.Set("User-Agent", lineSplit[index+1])
			}
		}
		return req, err
	})
}

var loadRequestBodyFromFile = func(filepath string) *bytes.Buffer {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		Log.Fatalf("Request body file not found: %s", filepath)
		return nil
	}
	Log.Println("file exists; processing...")
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		Log.Fatalf("Unable to read Request body file: %s", filepath)
		return nil
	}
	return bytes.NewBuffer(data)
}
