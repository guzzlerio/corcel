package yaml

import (
	"testing"

	yaml "github.com/ghodss/yaml"
	"github.com/guzzlerio/corcel/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlanBuilder(t *testing.T) {
	Convey("PlanBuilder", t, func() {

		Convey("builds a basic plan with defaults", func() {
			planBuilder := NewPlanBuilder()

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
`
			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("waitTime", func() {

			// Skipping this because gomega is unable to assert the panic if
			// the func returns a value.
			// Suggest switching to another assertion lib, but something like
			// github.com/stretchr/testify requires the *testing.T which
			// ginkgo does not take in
			SkipConvey("panics when waitTime does not parse", func() {
				planBuilder := NewPlanBuilder()

				planBuilder.SetWaitTime("talula")
			})
			Convey("builds a plan with waitTime", func() {
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

				So(string(planData), ShouldMatchYaml, expected)
			})
		})
		Convey("builds a plan with a name", func() {
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

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with random", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.SetRandom(true)

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
random: true
waitTime: 0s
workers: 1
duration: 0s
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with workers", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.SetWorkers(10)

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
random: false
waitTime: 0s
workers: 10
duration: 0s
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with iterations", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.SetIterations(10)

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
iterations: 10
random: false
waitTime: 0s
workers: 1
duration: 0s
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with duration", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.SetDuration("10s")

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
random: false
waitTime: 0s
workers: 1
duration: 10s
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with context", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.SetDuration("10s").
				WithContext(planBuilder.BuildContext().SetList("People", []core.ExecutionContext{
					{"name": "bob", "age": 52},
				}).Build())

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
random: false
waitTime: 0s
workers: 1
duration: 10s
context:
  lists:
    People:
    - age: 52
      name: bob
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with a before", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.WithName("Some Plan").
				AddBefore(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build())

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
name: Some Plan
random: false
waitTime: 0s
workers: 1
duration: 0s
before:
 - results:
     value:1: talula 123 bang bang
   type: DummyAction
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with a after", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.WithName("Some Plan").
				AddAfter(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build())

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
name: Some Plan
random: false
waitTime: 0s
workers: 1
duration: 0s
after:
 - results:
     value:1: talula 123 bang bang
   type: DummyAction
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("builds a plan with a before/after", func() {
			planBuilder := NewPlanBuilder()

			planBuilder.WithName("Some Plan").
				AddBefore(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
				AddAfter(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build())

			plan := planBuilder.Build()
			planData, _ := yaml.Marshal(&plan)

			var expected = `
name: Some Plan
random: false
waitTime: 0s
workers: 1
duration: 0s
before:
 - results:
     value:1: talula 123 bang bang
   type: DummyAction
after:
 - results:
     value:1: talula 123 bang bang
   type: DummyAction
`

			So(string(planData), ShouldMatchYaml, expected)
		})
		Convey("jobs", func() {
			Convey("builds a plan with a job with a name", func() {
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

				So(string(planData), ShouldMatchYaml, expected)
			})
			Convey("builds a plan with a job with a context", func() {
				planBuilder := NewPlanBuilder()

				planBuilder.CreateJob().
					WithContext(planBuilder.BuildContext().SetList("People", []core.ExecutionContext{
						{"name": "bob", "age": 52},
					}).Build())

				plan := planBuilder.Build()
				planData, _ := yaml.Marshal(&plan)

				var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- context:
    lists:
      People:
      - age: 52
        name: bob
`

				So(string(planData), ShouldMatchYaml, expected)
			})
			Convey("builds a plan with a job", func() {
				planBuilder := NewPlanBuilder()

				planBuilder.CreateJob("Some Job")

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

				So(string(planData), ShouldMatchYaml, expected)
			})
			Convey("builds a plan with a before/after job", func() {
				planBuilder := NewPlanBuilder()

				planBuilder.CreateJob().WithName("Some Job").
					AddBefore(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
					AddAfter(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build())

				plan := planBuilder.Build()
				planData, _ := yaml.Marshal(&plan)

				var expected = `
random: false
waitTime: 0s
workers: 1
duration: 0s
jobs:
- name: Some Job
  before:
  - results:
      value:1: talula 123 bang bang
    type: DummyAction
  after:
  - results:
      value:1: talula 123 bang bang
    type: DummyAction
`

				So(string(planData), ShouldMatchYaml, expected)
			})
			Convey("steps", func() {
				Convey("builds a plan with a step", func() {
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

					So(string(planData), ShouldMatchYaml, expected)
				})
				Convey("builds a plan with a before/after job", func() {
					planBuilder := NewPlanBuilder()

					planBuilder.CreateJob().WithName("Some Job").
						CreateStep().WithName("Some Step").
						AddBefore(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						AddAfter(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build())

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
    before:
    - results:
        value:1: talula 123 bang bang
      type: DummyAction
    after:
    - results:
        value:1: talula 123 bang bang
      type: DummyAction
`
					So(string(planData), ShouldMatchYaml, expected)
				})
				Convey("builds a plan with an action", func() {
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

					So(string(planData), ShouldMatchYaml, expected)
				})
				Convey("builds a plan with an assertion", func() {
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

					So(string(planData), ShouldMatchYaml, expected)
				})
				Convey("builds a plan with an extractor", func() {
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

					So(string(planData), ShouldMatchYaml, expected)
				})
			})
		})
	})
}
