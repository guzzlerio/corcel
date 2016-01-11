package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/logger"
)

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}
//RequestAdapter ...
type RequestAdapter struct {
	Handlers map[string]RequestConfigHandler
}

//NewRequestAdapter ...
func NewRequestAdapter() RequestAdapter {
	return RequestAdapter{
		Handlers: map[string]RequestConfigHandler{
			"-X": RequestConfigHandler(HandlerForMethod),
			"-H": RequestConfigHandler(HandlerForHeader),
			"-d": RequestConfigHandler(HandlerForData),
			"-A": RequestConfigHandler(HandlerForUserAgent),
		},
	}
}

//RequestConfigHandler ...
type RequestConfigHandler func(options []string, index int, req *http.Request) (*http.Request, error)

//HandlerForMethod ...
func HandlerForMethod(options []string, index int, req *http.Request) (*http.Request, error) {
	req.Method = options[index+1]
	return req, nil
}

//HandlerForHeader ...
func HandlerForHeader(options []string, index int, req *http.Request) (*http.Request, error) {
	value := strings.Trim(options[index+1], "\"")

	valueSplit := strings.Split(value, ":")
	req.Header.Set(strings.TrimSpace(valueSplit[0]), strings.TrimSpace(valueSplit[1]))
	return req, nil
}

//HandlerForData ...
func HandlerForData(options []string, index int, req *http.Request) (outReq *http.Request, err error) {
	rawBody := options[index+1]

	if strings.ToLower(req.Method) == "get" {
		req.URL.RawQuery = options[index+1]
		outReq = req
	} else {
		var body *bytes.Buffer
		bodyBytes := []byte(rawBody)
		if strings.HasPrefix(rawBody, "@") {
			body = loadRequestBodyFromFile(string(bytes.TrimLeft(bodyBytes, "@")))
		} else {
			logger.Log.Println("body from request")
			body = bytes.NewBuffer(bodyBytes)
		}
		outReq, err = http.NewRequest(req.Method, req.URL.String(), body)
	}
	return
}

//HandlerForUserAgent ...
func HandlerForUserAgent(options []string, index int, req *http.Request) (*http.Request, error) {
	req.Header.Set("User-Agent", options[index+1])
	return req, nil
}

//Create ...
func (instance RequestAdapter) Create(line string) RequestFunc {
	return RequestFunc(func() (*http.Request, error) {
		lexer := NewCommandLineLexer()
		lineSplit := lexer.Lex(line)
		req, err := http.NewRequest("GET", lineSplit[0], nil)
		if err != nil {
			return nil, err
		}
		for index := range lineSplit {
			arg := lineSplit[index]
			for key, handler := range instance.Handlers {
				if key == arg {
					req, err = handler(lineSplit, index, req)
				}
			}
		}
		return req, err
	})
}

var loadRequestBodyFromFile = func(filepath string) *bytes.Buffer {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		logger.Log.Fatalf("Request body file not found: %s\n", filepath)
		return nil
	}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Log.Fatalf("Unable to read Request body file: %s\n", filepath)
		return nil
	}
	return bytes.NewBuffer(data)
}
