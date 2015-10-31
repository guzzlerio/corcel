package main

import (
	"bytes"
	"net/http"
	"strings"
)

type RequestAdapter struct {}

func NewRequestAdapter() RequestAdapter {
	return RequestAdapter{}
}

func (instance RequestAdapter) Create(line string) RequestFunc {
	return RequestFunc(func() (*http.Request,error) {
        commandLineLexer := NewCommandLineLexer()
		lineSplit := commandLineLexer.Lex(line)
		req, err := http.NewRequest("GET", lineSplit[0], nil)
        if err != nil{
            return nil, err
        }
		for index, _ := range lineSplit {
			if lineSplit[index] == "-X" {
				req.Method = lineSplit[index+1]
			}
			if lineSplit[index] == "-H" {
				value := strings.Trim(lineSplit[index+1], "\"")

				valueSplit := strings.Split(value, ":")
				req.Header.Set(strings.TrimSpace(valueSplit[0]), strings.TrimSpace(valueSplit[1]))
			}
			if lineSplit[index] == "-d" {
				if strings.ToLower(req.Method) == "get" {
					req.URL.RawQuery = lineSplit[index+1]
				} else {
					body := bytes.NewBuffer([]byte(lineSplit[index+1]))
					req, err = http.NewRequest(req.Method, req.URL.String(), body)
				}
			}
		}
		return req, err
	})
}
