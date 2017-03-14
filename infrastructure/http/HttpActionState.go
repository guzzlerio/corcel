package http

import (
	"net/http"
)

type HttpActionState struct {
	URL     string
	Method  string
	Body    string
	Headers http.Header
}
