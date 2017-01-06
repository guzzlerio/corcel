package yaml

import (
	"io/ioutil"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("PlanBuilder", func() {

	ginkgo.It("Does something", func() {

		planBuilder := NewPlanBuilder()

		planBuilder.WithName("Some Plan").
			CreateJob().WithName("Some Job").
			CreateStep().WithName("Some Step").
			ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
			WithExtractor(planBuilder.RegexExtractor().Name("regex:match:1").Key("value:1").Match("\\d+").Build()).
			WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

		file, _ := planBuilder.Build()
		planData, _ := ioutil.ReadFile(file.Name())

		var expected = `duration: 0s
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
name: Some Plan
random: false
waitTime: 0s
workers: 1
`

		Expect(string(planData)).To(Equal(expected))
	})
})
