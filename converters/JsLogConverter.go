package converters

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/robertkrimen/otto"
)

// JsLogConverter
type JsLogConverter struct {
	baseUrl *url.URL
	scanner *bufio.Scanner
	fields  LogFields
	vm      *otto.Otto
}

// NewW3cExtConverter ...
func NewJsLogConverter(js string, baseUrl *url.URL, input io.Reader) *JsLogConverter {
	scanner := bufio.NewScanner(input)
	//TODO allow different scanner Split options see https://golang.org/pkg/bufio/index.html#Scanner.Split
	if baseUrl == nil {
		panic(fmt.Errorf("Base URL is required", nil))
	}
	vm := otto.New()
	if _, err := vm.Run(js); err != nil {
		panic(err)
	}

	return &JsLogConverter{
		baseUrl: baseUrl,
		scanner: scanner,
		vm:      vm,
		fields:  []string{},
	}
}

func (i *JsLogConverter) Convert() (*yaml.ExecutionPlan, error) {
	planBuilder := yaml.NewPlanBuilder()
	jobBuilder := planBuilder.CreateJob()
	for i.scanner.Scan() {
		line := i.scanner.Text()

		//TODO might not need this
		if isFields, _ := i.vm.Call("isFieldDefinition", nil, line); isFields == otto.TrueValue() {
			ottoLine, _ := i.vm.Call("parseLine", nil, line)
			if ottoFields, _ := ottoLine.Object().Get("fields"); ottoFields.IsDefined() {
				fields, _ := ottoFields.Export()
				i.fields = fields.([]string)
				continue
			}
		}

		entry := LogEntry{}
		ottoLine, err := i.vm.Call("parseLine", nil, line, i.fields)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(ottoLine.String()), &entry)

		if err := i.failsMinRequiredFields(entry); err != nil {
			fmt.Println(err)
			continue
		}

		jobBuilder.
			CreateStep().
			ToExecuteAction(planBuilder.HTTPAction().Method(entry.Request.Method).URL(i.buildURL(entry)).Build()).
			WithAssertion(planBuilder.ExactAssertion("response:status", entry.Response.Status))
	}
	if err := i.scanner.Err(); err != nil {
		return nil, err
	}
	plan := planBuilder.Build()
	plan.Name = fmt.Sprintf("Log file replay for: %v", i.baseUrl.String())
	return plan, nil
}

func (i *JsLogConverter) failsMinRequiredFields(entry LogEntry) error {
	if entry.Request.Method == "POST" && entry.Request.Payload == "" {
		return fmt.Errorf("POST method found in log entry but no payload")
	}
	if entry.Request.Host == "" && i.baseUrl == nil {
		return fmt.Errorf("No Host found in log entry and no base URL provided")
	}
	if entry.Request.Method == "" {
		return fmt.Errorf("No method found in log entry")
	}
	if entry.Response.Status < 100 {
		return fmt.Errorf("No expected response status found in log entry")
	}
	return nil
}

func (i *JsLogConverter) buildURL(entry LogEntry) string {
	u, _ := url.Parse(i.baseUrl.String())
	u.Path = entry.Request.Path
	u.RawQuery = entry.Request.Query
	return u.String()
}
