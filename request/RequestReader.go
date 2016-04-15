package request

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

//Func ...
type Func func() (*http.Request, error)

//Reader ...
type Reader struct {
	Requests []Func
}

//NewRequestReader ...
func NewRequestReader(filePath string) *Reader {
	file, err := os.Open(filePath)
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closes the file")
		}
	}()
	check(err)
	requests := []Func{}
	requestAdapter := NewRequestAdapter()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		requestFunc := requestAdapter.Create(line)
		requests = append(requests, requestFunc)
		check(err)
	}

	return &Reader{
		Requests: requests,
	}
}

//Size ...
func (instance *Reader) Size() int {
	return len(instance.Requests)
}

//Read ...
func (instance *Reader) Read(index int) Func {
	return instance.Requests[index]
}

//Lexer ...
type Lexer interface {
	Lex(args string) []string
}
