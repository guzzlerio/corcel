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
func NewJsLogConverter(js string, baseUrl string, input io.Reader) *JsLogConverter {
	scanner := bufio.NewScanner(input)
	//TODO allow different scanner Split options see https://golang.org/pkg/bufio/index.html#Scanner.Split
	//TODO test for error
	u, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}
	vm := otto.New()
	if _, err := vm.Run(js); err != nil {
		panic(err)
	}

	return &JsLogConverter{
		baseUrl: u,
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
		ottoLine, _ := i.vm.Call("parseLine", nil, line, i.fields)
		json.Unmarshal([]byte(ottoLine.String()), &entry)

		if i.failsMinRequiredFields(entry) {
			panic(fmt.Errorf("Insufficient populated fields to convert: %+v", entry))
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

func (i *JsLogConverter) failsMinRequiredFields(entry LogEntry) bool {
	return false
}

func (i *JsLogConverter) buildURL(entry LogEntry) string {
	u, _ := url.Parse(i.baseUrl.String())
	u.Path = entry.Request.Path
	u.RawQuery = entry.Request.Query
	return u.String()
}
