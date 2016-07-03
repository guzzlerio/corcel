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
	ExecutionExtractorParsers map[string]core.ExecutionExtractorParser
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
	executionPlan.Context = yamlExecutionPlan.Context
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
		job := executionPlan.CreateJob()
		job.Name = yamlJob.Name

		for _, yamlStep := range yamlJob.Steps {
			step := job.CreateStep()
			step.Name = yamlStep.Name

			if yamlStep.Action["type"] != nil {
				actionType := yamlStep.Action["type"].(string)

				if parser := instance.ExecutionActionParsers[actionType]; parser != nil {
					step.Action = parser.Parse(yamlStep.Action)
				} else {
					panic(fmt.Sprintf("No parser configured for action %s", actionType))
				}
			}

			for _, yamlAssertion := range yamlStep.Assertions {
				assertionType := yamlAssertion["type"].(string)
				if parser := instance.ExecutionAssertionParsers[assertionType]; parser != nil {
					step.Assertions = append(step.Assertions, parser.Parse(yamlAssertion))
				} else {
					panic(fmt.Sprintf("No parser configured for assertion %s", assertionType))
				}
			}

			for _, yamlExtractor := range yamlStep.Extractors {
				extractorType := yamlExtractor["type"].(string)
				if parser := instance.ExecutionExtractorParsers[extractorType]; parser != nil {
					step.Extractors = append(step.Extractors, parser.Parse(yamlExtractor))
				} else {
					panic(fmt.Sprintf("No parser configured for extractor %s", extractorType))
				}
			}

			job = job.AddStep(step)
		}

		executionPlan = executionPlan.AddJob(job)
	}

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

//AddExtractorParser ...
func (instance *ExecutionPlanParser) AddExtractorParser(assertionType string, parser core.ExecutionExtractorParser) {
	if instance.ExecutionExtractorParsers == nil {
		instance.ExecutionExtractorParsers = map[string]core.ExecutionExtractorParser{}
	}
	instance.ExecutionExtractorParsers[assertionType] = parser
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

	//This can be refactored so that the Key method is invoked inside the AddActionParser
	for _, extractorParser := range registry.ExtractorParsers {
		parser.AddExtractorParser(extractorParser.Key(), extractorParser)
	}
	return parser
}
