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

		var expected = `random: false
workers: 1
waitTime: 0s
duration: 0s
name: Some Plan
jobs:
- name: Some Job
  steps:
  - name: Some Step
    action:
      results:
        value:1: talula 123 bang bang
      type: DummyAction
    extractors:
    - key: value:1
      match: \d+
      name: regex:match:1
      type: RegexExtractor
    assertions:
    - expected: "123"
      key: regex:match:1
      type: ExactAssertion
`

		Expect(string(planData)).To(Equal(expected))
	})
})
