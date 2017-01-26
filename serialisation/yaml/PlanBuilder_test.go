package yaml

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ghodss/yaml"
)

var _ = Describe("PlanBuilder", func() {

	It("builds a basic plan with defaults", func() {
		planBuilder := NewPlanBuilder()

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with waitTime", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.SetWaitTime("10s")

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
random: false
waitTime: 10s
workers: 1
duration: 0s
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with a name", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.WithName("Some Plan")

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
name: Some Plan
random: false
waitTime: 0s
workers: 1
duration: 0s
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with a job", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.CreateJob().WithName("Some Job")

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- name: Some Job
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with a step", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.
			CreateJob().WithName("Some Job").
			CreateStep().WithName("Some Step")

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- name: Some Job
  steps:
  - name: Some Step
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with an action", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.
			CreateJob().WithName("Some Job").
			CreateStep().WithName("Some Step").
			ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build())

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- name: Some Job
  steps:
  - action:
      results:
        value:1: talula 123 bang bang
      type: DummyAction
    name: Some Step
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with an assertion", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.
			CreateJob().WithName("Some Job").
			CreateStep().WithName("Some Step").
			ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
			WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- name: Some Job
  steps:
  - action:
      results:
        value:1: talula 123 bang bang
      type: DummyAction
    assertions:
    - expected: "123"
      key: regex:match:1
      type: ExactAssertion
    name: Some Step
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
	It("builds a plan with an extractor", func() {
		planBuilder := NewPlanBuilder()

		planBuilder.WithName("Some Plan").
			CreateJob().WithName("Some Job").
			CreateStep().WithName("Some Step").
			ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
			WithExtractor(planBuilder.RegexExtractor().Name("regex:match:1").Key("value:1").Match("\\d+").Build()).
			WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

		plan := planBuilder.Build()
		planData, _ := yaml.Marshal(&plan)

		var expected = `
name: Some Plan
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- name: Some Job
  steps:
  - action:
      results:
        value:1: talula 123 bang bang
      type: DummyAction
    assertions:
    - expected: "123"
      key: regex:match:1
      type: ExactAssertion
    extractors:
    - key: value:1
      match: \d+
      name: regex:match:1
      type: RegexExtractor
    name: Some Step
`

		Expect(string(planData)).To(MatchYAML(expected))
	})
})
