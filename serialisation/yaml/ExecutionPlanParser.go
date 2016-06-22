package yaml

import (
	"fmt"
	"time"

	"ci.guzzler.io/guzzler/corcel/core"
	"gopkg.in/yaml.v2"
)

//ExecutionPlanParser ...
type ExecutionPlanParser struct {
	ExecutionActionParsers    map[string]core.ExecutionActionParser
	ExecutionAssertionParsers map[string]core.ExecutionAssertionParser
}

//Parse ...
func (instance *ExecutionPlanParser) Parse(data string) (core.Plan, error) {
	var executionPlan core.Plan
	var yamlExecutionPlan ExecutionPlan

	err := yaml.Unmarshal([]byte(data), &yamlExecutionPlan)

	if err != nil {
		return core.Plan{}, err
	}

	executionPlan.Name = yamlExecutionPlan.Name
	executionPlan.WaitTime, err = time.ParseDuration(yamlExecutionPlan.WaitTime)
	if err != nil {
		executionPlan.WaitTime = time.Duration(0)
	}

	executionPlan.Duration, err = time.ParseDuration(yamlExecutionPlan.Duration)
	if err != nil {
		executionPlan.Duration = time.Duration(0)
	}

	executionPlan.Random = yamlExecutionPlan.Random

	executionPlan.Workers = yamlExecutionPlan.Workers

	for _, yamlJob := range yamlExecutionPlan.Jobs {
		job := core.Job{
			Name: yamlJob.Name,
		}

		for _, yamlStep := range yamlJob.Steps {
			step := core.Step{
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
					panic(fmt.Sprintf("No parser configured for assertion %s", assertionType))
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
func (instance *ExecutionPlanParser) AddActionParser(actionType string, parser core.ExecutionActionParser) {
	if instance.ExecutionActionParsers == nil {
		instance.ExecutionActionParsers = map[string]core.ExecutionActionParser{}
	}
	instance.ExecutionActionParsers[actionType] = parser
}

//AddAssertionParser ...
func (instance *ExecutionPlanParser) AddAssertionParser(assertionType string, parser core.ExecutionAssertionParser) {
	if instance.ExecutionAssertionParsers == nil {
		instance.ExecutionAssertionParsers = map[string]core.ExecutionAssertionParser{}
	}
	instance.ExecutionAssertionParsers[assertionType] = parser
}

//CreateExecutionPlanParser ...
func CreateExecutionPlanParser(registry core.Registry) *ExecutionPlanParser {
	parser := &ExecutionPlanParser{}

	//This can be refactored so that the Key method is invoked inside the AddActionParser
	for _, actionParser := range registry.ActionParsers {
		parser.AddActionParser(actionParser.Key(), actionParser)
	}

	//This can be refactored so that the Key method is invoked inside the AddActionParser
	for _, assertionParser := range registry.AssertionParsers {
		parser.AddAssertionParser(assertionParser.Key(), assertionParser)
	}
	return parser
}