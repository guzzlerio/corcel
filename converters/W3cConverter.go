package converters

import (
	"bufio"
	"fmt"
	"io"

	"github.com/guzzlerio/corcel/serialisation/yaml"
)

// W3cConverter
type W3cExtConverter struct {
	baseUrl string
	scanner *bufio.Scanner
}

// NewW3cExtConverter ...
func NewW3cExtConverter(baseUrl string, input io.Reader) *W3cExtConverter {
	scanner := bufio.NewScanner(input)
	return &W3cExtConverter{
		baseUrl: baseUrl,
		scanner: scanner,
	}
}

func (i *W3cExtConverter) Convert() (*yaml.ExecutionPlan, error) {
	 
	plan := yaml.ExecutionPlan{
		Jobs: 
	}
	for i.scanner.Scan() {
		fmt.Println(i.scanner.Text())
	}
	if err := i.scanner.Err(); err != nil {
		return nil, err
	}

	return &plan, nil
}
