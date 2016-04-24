package processor

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"

	"ci.guzzler.io/guzzler/corcel/logger"

	"gopkg.in/yaml.v2"
)

type ExecutionResultProcessor interface {
	Process(result ExecutionResult, registry metrics.Registry)
}

func NewGeneralExecutionResultProcessor() GeneralExecutionResultProcessor {
	return GeneralExecutionResultProcessor{}
}

type GeneralExecutionResultProcessor struct {
}

func (instance GeneralExecutionResultProcessor) Process(result ExecutionResult, registry metrics.Registry) {
	obj := result["action:duration"]
	timer := metrics.GetOrRegisterTimer("action:duration", registry)
	timer.Update(obj.(time.Duration))

	throughput := metrics.GetOrRegisterMeter("action:throughput", registry)
	throughput.Mark(1)

	errors := metrics.GetOrRegisterMeter("action:error", registry)
	if result["action:error"] != nil {
		var errorString string

		switch t := result["action:error"].(type) {
		case error:
			errorString = t.Error()
		case string:
			errorString = t
		}
		if !strings.Contains(errorString, "net/http: request canceled") {
			errors.Mark(1)
		}
	}

	if result["action:bytes:sent"] != nil {
		bytesSentValue := int64(result["action:bytes:sent"].(int))

		bytesSentCounter := metrics.GetOrRegisterCounter("counter:action:bytes:sent", registry)
		bytesSentCounter.Inc(bytesSentValue)

		bytesSent := metrics.GetOrRegisterHistogram("histogram:action:bytes:sent", registry, metrics.NewUniformSample(100))
		bytesSent.Update(bytesSentValue)
	}

	if result["action:bytes:received"] != nil {
		bytesReceivedValue := int64(result["action:bytes:received"].(int))

		bytesReceivedCounter := metrics.GetOrRegisterCounter("counter:action:bytes:received", registry)
		bytesReceivedCounter.Inc(bytesReceivedValue)

		bytesReceived := metrics.GetOrRegisterHistogram("histogram:action:bytes:received", registry, metrics.NewUniformSample(100))
		bytesReceived.Update(int64(result["action:bytes:received"].(int)))
	}
}

func NewHTTPExecutionResultProcessor() HTTPExecutionResultProcessor {
	return HTTPExecutionResultProcessor{}
}

type HTTPExecutionResultProcessor struct {
}

func (instance HTTPExecutionResultProcessor) Process(result ExecutionResult, registry metrics.Registry) {
	for key, value := range result {
		switch key {
		case "http:request:error":
			meter := metrics.GetOrRegisterMeter("http:request:error", registry)
			meter.Mark(1)

			url := result["http:request:url"]

			byUrlRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byUrl:%s:", url))
			byUrlmeter := metrics.GetOrRegisterMeter("http:request:error", byUrlRegistry)
			byUrlmeter.Mark(1)

		case "http:response:error":
			meter := metrics.GetOrRegisterMeter("http:response:error", registry)
			meter.Mark(1)

			url := result["http:request:url"]
			byUrlRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byUrl:%s:", url))
			byUrlmeter := metrics.GetOrRegisterMeter("http:response:error", byUrlRegistry)
			byUrlmeter.Mark(1)

		case "http:request:bytes":
			url := result["http:request:url"]
			byUrlRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byUrl:%s:", url))
			byUrlHistogram := metrics.GetOrRegisterHistogram("http:request:bytes", byUrlRegistry, metrics.NewUniformSample(100))
			byUrlHistogram.Update(int64(value.(int)))

		case "http:response:bytes":
			url := result["http:request:url"]
			byUrlRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byUrl:%s:", url))
			byUrlHistogram := metrics.GetOrRegisterHistogram("http:response:bytes", byUrlRegistry, metrics.NewUniformSample(100))
			byUrlHistogram.Update(int64(value.(int)))

		case "http:response:status":
			statusCode := value.(int)
			url := result["http:request:url"]
			counter := metrics.GetOrRegisterCounter(fmt.Sprintf("http:response:status:%d", statusCode), registry)
			counter.Inc(1)

			byUrlRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byUrl:%s:", url))
			byUrlCounter := metrics.GetOrRegisterCounter(fmt.Sprintf("http:response:status:%d", statusCode), byUrlRegistry)
			byUrlCounter.Inc(1)

			obj := result["action:duration"]
			timer := metrics.GetOrRegisterTimer(fmt.Sprintf("http:response:status:%d:duration", statusCode), registry)
			timer.Update(obj.(time.Duration))

			byUrlTimer := metrics.GetOrRegisterTimer(fmt.Sprintf("http:response:status:%d:duration", statusCode), byUrlRegistry)
			byUrlTimer.Update(obj.(time.Duration))
		}
	}
}

//HTTPRequestExecutionAction ...
type HTTPRequestExecutionAction struct {
	Client  *http.Client
	URL     string
	Method  string
	Headers http.Header
}

func (instance *HTTPRequestExecutionAction) initialize() {
	instance.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}
}

//Execute ...
func (instance *HTTPRequestExecutionAction) Execute(cancellation chan struct{}) ExecutionResult {
	if instance.Client == nil {
		instance.initialize()
	}

	result := ExecutionResult{}

	req, err := http.NewRequest(instance.Method, instance.URL, nil)
	req.Cancel = cancellation
	//This should be a configuration item.  It allows the client to work
	//in a way similar to a server which does not support HTTP KeepAlive
	//After each request the client channel is closed.  When set to true
	//the performance overhead is large in terms of Network IO throughput

	//req.Close = true

	if err != nil {
		panic(err)
		fmt.Println(fmt.Sprintf("HTTP Err %v", err))
		result["action:error"] = err
		return result
	}

	req.Header = instance.Headers
	response, err := instance.Client.Do(req)
	if err != nil {
		result["action:error"] = err
		return result
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Log.Warnf("Error closing response Body %v", err)
		}
	}()

	requestBytes, _ := httputil.DumpRequest(req, true)
	responseBytes, _ := httputil.DumpResponse(response, true)

	if response.StatusCode >= 500 {
		result["action:error"] = fmt.Sprintf("Server Error %d", response.StatusCode)
	}

	result["http:request:url"] = req.URL.String()
	result["action:bytes:sent"] = len(requestBytes)
	result["action:bytes:received"] = len(responseBytes)
	result["http:request:headers"] = req.Header
	result["http:response:status"] = response.StatusCode

	return result
}

//YamlHTTPRequestParser ...
type YamlHTTPRequestParser struct{}

//Parse ...
func (instance YamlHTTPRequestParser) Parse(input map[string]interface{}) Action {
	action := HTTPRequestExecutionAction{
		URL:     input["url"].(string),
		Method:  input["method"].(string),
		Headers: http.Header{},
	}
	for key, value := range input["httpHeaders"].(map[interface{}]interface{}) {
		action.Headers.Set(key.(string), value.(string))
	}
	return &action
}

//Key ...
func (instance YamlHTTPRequestParser) Key() string {
	return "HttpRequest"
}

//YamlExactAssertionParser ...
type YamlExactAssertionParser struct{}

//Parse ...
func (instance YamlExactAssertionParser) Parse(input map[string]interface{}) Assertion {
	return &ExactAssertion{
		Key:      input["key"].(string),
		Expected: input["expected"].(int),
	}
}

//Key ...
func (instance YamlExactAssertionParser) Key() string {
	return "ExactAssertion"
}

//ExactAssertion ...
type ExactAssertion struct {
	Key      string
	Expected interface{}
}

//ResultKey ...
func (instance *ExactAssertion) ResultKey() string {
	return instance.Key + ":assert:exactmatch"
}

//Assert ...
func (instance *ExactAssertion) Assert(executionResult ExecutionResult) AssertionResult {
	actual := executionResult[instance.Key]

	result := map[string]interface{}{
		"expected": instance.Expected,
		"actual":   actual,
	}
	if actual == instance.Expected {
		result["result"] = "pass"
	} else {
		result["result"] = "fail"
		result["message"] = fmt.Sprintf("FAIL: %v does not match %v", actual, instance.Expected)
	}
	return result
}

//YamlExecutionStep ...
type YamlExecutionStep struct {
	Name       string                   `yaml:"name"`
	Action     map[string]interface{}   `yaml:"action"`
	Extract    map[string]string        `yaml:"extract"`
	Assertions []map[string]interface{} `yaml:"assertions"`
}

//YamlExecutionJob ...
type YamlExecutionJob struct {
	Name  string              `yaml:"name"`
	Steps []YamlExecutionStep `yaml:"steps"`
}

//YamlExecutionPlan ...
type YamlExecutionPlan struct {
	Name     string             `yaml:"name"`
	Workers  int                `yaml:"workers"`
	WaitTime string             `yaml:"waitTime"`
	Jobs     []YamlExecutionJob `yaml:"jobs"`
}

//YamlExecutionActionParser ...
type YamlExecutionActionParser interface {
	Parse(input map[string]interface{}) Action
	Key() string
}

//YamlExecutionAssertionParser ...
type YamlExecutionAssertionParser interface {
	Parse(input map[string]interface{}) Assertion
	Key() string
}

//ExecutionResult ...
type ExecutionResult map[string]interface{}

//AssertionResult ...
type AssertionResult map[string]interface{}

//Action ...
type Action interface {
	Execute(cancellation chan struct{}) ExecutionResult
}

//Assertion ...
type Assertion interface {
	ResultKey() string
	Assert(ExecutionResult) AssertionResult
}

//Step ...
type Step struct {
	Name       string
	Action     Action
	Assertions []Assertion
}

//Job ...
type Job struct {
	Name  string
	Steps []Step
}

//Plan ...
type Plan struct {
	Name     string
	Workers  int
	WaitTime time.Duration
	Jobs     []Job
}

//ExecutionPlanParser ...
type ExecutionPlanParser struct {
	ExecutionActionParsers    map[string]YamlExecutionActionParser
	ExecutionAssertionParsers map[string]YamlExecutionAssertionParser
}

//Parse ...
func (instance *ExecutionPlanParser) Parse(data string) (Plan, error) {
	var executionPlan Plan
	var yamlExecutionPlan YamlExecutionPlan

	err := yaml.Unmarshal([]byte(data), &yamlExecutionPlan)

	if err != nil {
		return Plan{}, err
	}

	executionPlan.Name = yamlExecutionPlan.Name
	executionPlan.WaitTime, err = time.ParseDuration(yamlExecutionPlan.WaitTime)
	if err != nil {
		executionPlan.WaitTime = time.Duration(0)
	}
	executionPlan.Workers = yamlExecutionPlan.Workers

	for _, yamlJob := range yamlExecutionPlan.Jobs {
		job := Job{
			Name: yamlJob.Name,
		}

		for _, yamlStep := range yamlJob.Steps {
			step := Step{
				Name: yamlStep.Name,
			}
			actionType := yamlStep.Action["type"].(string)

			if parser := instance.ExecutionActionParsers[actionType]; parser != nil {
				step.Action = parser.Parse(yamlStep.Action)
			} else {
				panic(fmt.Sprintf("No parser configured for action %s", actionType))
			}
			for _, yamlAssertion := range yamlStep.Assertions {
				assertionType := yamlAssertion["type"].(string)
				if parser := instance.ExecutionAssertionParsers[assertionType]; parser != nil {
					step.Assertions = append(step.Assertions, parser.Parse(yamlAssertion))
				} else {
					panic(fmt.Sprintf("No parser configured for action %s", actionType))
				}
			}

			job.Steps = append(job.Steps, step)
		}

		executionPlan.Jobs = append(executionPlan.Jobs, job)
	}

	//We have an execution plan

	//Now we need to execute it.

	return executionPlan, nil
}

//AddActionParser ...
func (instance *ExecutionPlanParser) AddActionParser(actionType string, parser YamlExecutionActionParser) {
	if instance.ExecutionActionParsers == nil {
		instance.ExecutionActionParsers = map[string]YamlExecutionActionParser{}
	}
	instance.ExecutionActionParsers[actionType] = parser
}

//AddAssertionParser ...
func (instance *ExecutionPlanParser) AddAssertionParser(assertionType string, parser YamlExecutionAssertionParser) {
	if instance.ExecutionAssertionParsers == nil {
		instance.ExecutionAssertionParsers = map[string]YamlExecutionAssertionParser{}
	}
	instance.ExecutionAssertionParsers[assertionType] = parser
}

//CreateExecutionPlanParser ...
func CreateExecutionPlanParser() *ExecutionPlanParser {
	parser := &ExecutionPlanParser{}
	actionParsers := []YamlExecutionActionParser{YamlHTTPRequestParser{}}
	assertionParsers := []YamlExecutionAssertionParser{YamlExactAssertionParser{}}

	//This can be refactored so that the Key method is invoked inside the AddActionParser
	for _, actionParser := range actionParsers {
		parser.AddActionParser(actionParser.Key(), actionParser)
	}

	//This can be refactored so that the Key method is invoked inside the AddActionParser
	for _, assertionParser := range assertionParsers {
		parser.AddAssertionParser(assertionParser.Key(), assertionParser)
	}
	return parser
}
