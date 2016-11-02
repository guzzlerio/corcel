package main

import "ci.guzzler.io/guzzler/corcel/infrastructure/http"

//GetPathRequest ...
func GetHTTPRequestAction(url string) map[string]interface{} {
	return map[string]interface{}{
		"type":           "HttpRequest",
		"requestTimeout": 150,
		"method":         "GET",
		"url":            TestServer.CreateURL(url),
		"httpHeaders": map[string]string{
			"Content-Type": "application/json",
		},
	}
}

//HTTPStatusExactAssertion ...
func HTTPStatusExactAssertion(code int) map[string]interface{} {
	return map[string]interface{}{
		"type":     "ExactAssertion",
		"key":      http.ResponseStatusUrn.String(),
		"expected": code,
	}
}
