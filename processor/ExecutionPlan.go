package processor

import (
	"fmt"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"

	"ci.guzzler.io/guzzler/corcel/core"
	"ci.guzzler.io/guzzler/corcel/infrastructure/http"

	"gopkg.in/yaml.v2"
)

//ExecutionResultProcessor ...
type ExecutionResultProcessor interface {
	Process(result core.ExecutionResult, registry metrics.Registry)
}

//NewGeneralExecutionResultProcessor ...
func NewGeneralExecutionResultProcessor() GeneralExecutionResultProcessor {
	return GeneralExecutionResultProcessor{}
}

//GeneralExecutionResultProcessor ...
type GeneralExecutionResultProcessor struct {
}

//Process ...
func (instance GeneralExecutionResultProcessor) Process(result core.ExecutionResult, registry metrics.Registry) {
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
func (instance *ExactAssertion) Assert(executionResult core.ExecutionResult) core.AssertionResult {
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
	Random   bool               `yaml:"random"`
	Workers  int                `yaml:"workers"`
	WaitTime string             `yaml:"waitTime"`
	Duration string             `yaml:"duration"`
	Name     string             `yaml:"name"`
	Jobs     []YamlExecutionJob `yaml:"jobs"`
}

//YamlExecutionActionParser ...
type YamlExecutionActionParser interface {
	Parse(input map[string]interface{}) core.Action
	Key() string
}

//YamlExecutionAssertionParser ...
type YamlExecutionAssertionParser interface {
	Parse(input map[string]interface{}) Assertion
	Key() string
}

//Assertion ...
type Assertion interface {
	ResultKey() string
	Assert(core.ExecutionResult) core.AssertionResult
}

//Step ...
type Step struct {
	Name       string
	Action     core.Action
	Assertions []Assertion
}

//Job ...
type Job struct {
	Name  string
	Steps []Step
}

//Plan ...
type Plan struct {
	Random   bool
	Workers  int
	Name     string
	WaitTime time.Duration
	Duration time.Duration
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

	executionPlan.Duration, err = time.ParseDuration(yamlExecutionPlan.Duration)
	fmt.Println(fmt.Sprintf("THE Duration %v", yamlExecutionPlan.Duration))
	if err != nil {
		executionPlan.Duration = time.Duration(0)
	}

	executionPlan.Random = yamlExecutionPlan.Random

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
	actionParsers := []YamlExecutionActionParser{http.YamlHTTPRequestParser{}}
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
