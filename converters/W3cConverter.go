package converters

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/guzzlerio/corcel/serialisation/yaml"
)

type w3cFields []string

// W3cConverter
type W3cExtConverter struct {
	baseUrl *url.URL
	scanner *bufio.Scanner
	fields  w3cFields
}

type logEntry map[string]string

// NewW3cExtConverter ...
func NewW3cExtConverter(baseUrl string, input io.Reader) *W3cExtConverter {
	scanner := bufio.NewScanner(input)
	//TODO check for error
	u, _ := url.Parse(baseUrl)
	return &W3cExtConverter{
		baseUrl: u,
		scanner: scanner,
	}
}

func (i *W3cExtConverter) Convert() (*yaml.ExecutionPlan, error) {
	planBuilder := yaml.NewPlanBuilder()
	jobBuilder := planBuilder.CreateJob()
	for i.scanner.Scan() {
		line := i.scanner.Text()
		if i.isDirective(line) {
			if i.isFieldDefinition(line) {
				i.parseFields(line)
			}
			continue
		}
		entry := i.parseLine(line)
		if i.failsMinRequiredFields(entry) {
			panic(fmt.Errorf("Insufficient populated fields to convert: %+v", entry))
		}
		expectedStatus, _ := strconv.Atoi(entry["sc-status"])
		jobBuilder.
			CreateStep().
			ToExecuteAction(planBuilder.HTTPAction().Method(entry["cs-method"]).URL(i.buildURL(entry)).Build()).
			WithAssertion(planBuilder.ExactAssertion("response:status", expectedStatus))
	}
	if err := i.scanner.Err(); err != nil {
		return nil, err
	}
	plan := planBuilder.Build()
	if file, err := planBuilder.Write(plan); err == nil {
		defer func() {
			if fileErr := os.Remove(file.Name()); fileErr != nil {
				panic(fileErr)
			}
		}()
		dat, _ := ioutil.ReadFile(file.Name())

		fmt.Sprintln(string(dat))
	}

	return plan, nil
}

func (i *W3cExtConverter) isDirective(line string) bool {
	return strings.HasPrefix(line, "#")
}

func (i *W3cExtConverter) isFieldDefinition(line string) bool {
	return strings.HasPrefix(line, "#Fields: ")
}

func (i *W3cExtConverter) parseFields(line string) {
	i.fields = strings.Split(strings.TrimPrefix(line, "#Fields: "), " ")
}

func (i *W3cExtConverter) parseLine(line string) logEntry {
	a := strings.Split(line, " ")
	if len(a) != len(i.fields) {
		panic(fmt.Errorf("Log line entries does not match Field definition: %v - %v", len(a), len(i.fields)))
	}
	result := logEntry{}
	for i, v := range i.fields {
		result[v] = a[i]
	}
	return result
}

func (i *W3cExtConverter) failsMinRequiredFields(entry logEntry) bool {
	return false
}

func (i *W3cExtConverter) buildURL(entry logEntry) string {
	u, _ := url.Parse(i.baseUrl.String())
	u.Path = entry["cs-uri-stem"]
	if val, ok := entry["cs-uri-query"]; ok {
		if val != "-" {
			u.RawQuery = val
		}
	}
	return u.String()
}
