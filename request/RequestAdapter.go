package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
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

//Adapter ...
type Adapter struct {
	Handlers map[string]ConfigHandler
}

//NewRequestAdapter ...
func NewRequestAdapter() Adapter {
	return Adapter{
		Handlers: map[string]ConfigHandler{
			"-X": ConfigHandler(HandlerForMethod),
			"-H": ConfigHandler(HandlerForHeader),
			"-d": ConfigHandler(HandlerForData),
			"-A": ConfigHandler(HandlerForUserAgent),
		},
	}
}

//ConfigHandler ...
type ConfigHandler func(options []string, index int, req *http.Request) (*http.Request, error)

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
func (instance Adapter) Create(line string) Func {
	return Func(func() (*http.Request, error) {
		lexer := NewCommandLineLexer()
		lineSplit := lexer.Lex(line)
		requestURL := lineSplit[0]
		_, err := url.ParseRequestURI(requestURL)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("GET", requestURL, nil)
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
