package processor

import (
	//	"github.com/robertkrimen/otto"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

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
func (instance *HTTPRequestExecutionAction) Execute() (ExecutionResult, error) {
	if instance.Client == nil {
		instance.initialize()
	}

	req, _ := http.NewRequest(instance.Method, instance.URL, nil)
	/*
		if len(instance.Headers) > 0 {
			for key, value := range instance.Headers {
				req.Header.Set(key, value)
			}
		}
	*/
	req.Header = instance.Headers
	response, _ := instance.Client.Do(req)
	value := ExecutionResult{
		"http:request:headers": req.Header,
		"http:response:status": response.StatusCode,
	}
	return value, nil
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

//YamlExactAssertionParser ...
type YamlExactAssertionParser struct{}

//Parse ...
func (instance YamlExactAssertionParser) Parse(input map[string]interface{}) Assertion {
	return &ExactAssertion{
		Key:      input["key"].(string),
		Expected: input["expected"].(int),
	}
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
	WaitTime time.Duration      `yaml:"waitTime"`
	Jobs     []YamlExecutionJob `yaml:"jobs"`
}

//YamlExecutionActionParser ...
type YamlExecutionActionParser interface {
	Parse(input map[string]interface{}) Action
}

//YamlExecutionAssertionParser ...
type YamlExecutionAssertionParser interface {
	Parse(input map[string]interface{}) Assertion
}

//ExecutionResult ...
type ExecutionResult map[string]interface{}

//AssertionResult ...
type AssertionResult map[string]interface{}

//Action ...
type Action interface {
	Execute() (ExecutionResult, error)
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

//Execute ...
func (instance *Plan) Execute() error {
	//resultChannel := make(chan ExecutionResult)

	for _, job := range instance.Jobs {
		func(talula Job) {
			for _, step := range talula.Steps {
				executionResult, _ := step.Action.Execute()
				for _, assertion := range step.Assertions {
					assertionResult := assertion.Assert(executionResult)
					executionResult[assertion.ResultKey()] = assertionResult
				}
			}
		}(job)
	}

	return nil
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

	PrintYamlExecutionPlan(yamlExecutionPlan)

	if err != nil {
		return Plan{}, err
	}

	executionPlan.Name = yamlExecutionPlan.Name
	executionPlan.WaitTime = yamlExecutionPlan.WaitTime
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

//PrintYamlExecutionPlan ...
func PrintYamlExecutionPlan(plan YamlExecutionPlan) {
	fmt.Println(fmt.Sprintf("%v", plan.Name))
	fmt.Println(fmt.Sprintf("%v", plan.Jobs[0].Name))
}

/*
func main() {
	path := "./test-data.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panic("BOOOOOOM")
	}
	data, _ := ioutil.ReadFile(path)

	var parser ExecutionPlanParser
	parser.AddActionParser("HttpRequest", YamlHTTPRequestParser{})
	parser.AddAssertionParser("ExactAssertion", YamlExactAssertionParser{})

	executionPlan, err := parser.Parse(string(data))

	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	err = executionPlan.Execute()

	//Each job gets its own go routine
	//    In each job it iterates over the steps until it should stop
	//      - One iteration = One Job iterating every step
	//        - run 5 times that would be run every step in each job 5 times.
	//      - run for 10 seconds will stop any step when the time elpses
	//      - random makes one job randomise the order of execution for its steps
	//
	//      - TODO: support x number of iterations

	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

}
*/
