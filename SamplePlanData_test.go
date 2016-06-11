package main

//GetPathRequest ...
func GetPathRequest(url string) map[string]interface{} {
	return map[string]interface{}{
		"type":          "HttpRequest",
		"requesTimeout": 150,
		"method":        "GET",
		"url":           TestServer.CreateURL(url),
		"httpHeaders": map[string]string{
			"Content-Type": "application/json",
		},
	}
}

//HTTPStatusExactAssertion ...
func HTTPStatusExactAssertion(code int) map[string]interface{} {
	return map[string]interface{}{
		"type":     "ExactAssertion",
		"key":      "http:response:status",
		"expected": code,
	}
}
