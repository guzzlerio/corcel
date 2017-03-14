package main

import (
	"net/http"
)

func returnCode200(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.Header().Set("Content-Type", "text/plain")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", returnCode200)

	http.ListenAndServe(":1337", mux)
}
