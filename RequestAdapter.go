package main

import (
	"net/http"
	"strings"
)

type RequestAdapter struct {
}

func NewRequestAdapter() RequestAdapter {
	return RequestAdapter{}
}

func (instance RequestAdapter) Create(line string) (*http.Request, error) {
	line = strings.TrimSpace(line)
	lineSplit := strings.Split(line, " ")
	if len(lineSplit) > 1 {
		req, err := http.NewRequest("GET", lineSplit[0], nil)
		if lineSplit[1] == "-X" {
			req.Method = lineSplit[2]
		}
		return req, err
	} else {
		return http.NewRequest("GET", line, nil)
	}
}
