package main

import (
	"os"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/test"
	"github.com/guzzlerio/corcel/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanContexts", func() {

	contextsDebugPath := "/tmp/contexts"
	Describe("Lists", func() {

		BeforeEach(func() {
			if err := os.Remove(contextsDebugPath); err != nil {

			}
		})

		AfterEach(func() {
			if err := os.Remove(contextsDebugPath); err != nil {

			}
		})

		Context("Plan Scope", func() {

			It("Succeeds", func() {

				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					SetIterations(3).
					WithContext(planBuilder.BuildContext().SetList("People", []map[string]interface{}{
						{"name": "jill", "age": 35},
						{"name": "bob", "age": 52},
						{"name": "carol", "age": 24},
					}).Build()).
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().LogToFile(contextsDebugPath).Build())

				err := ExecutePlanBuilder(planBuilder)
				Expect(err).To(BeNil())

				contexts := test.GetExecutionContexts(contextsDebugPath)
				Expect(len(contexts)).To(Equal(3))
				Expect(contexts[0]["$People.name"]).To(Equal("jill"))
				Expect(contexts[1]["$People.name"]).To(Equal("bob"))
				Expect(contexts[2]["$People.name"]).To(Equal("carol"))
			})

		})
	})

	Context("Plan Scope", func() {
		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build()).
				WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("$value:3", "3")).
				WithAssertion(planBuilder.ExactAssertion("something", "1"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})
		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("$value:3", "3"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(3)))
		})
	})

	Context("Job Scope", func() {
		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build()).
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("$value:3", "3"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build())

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("$value:3", "3"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(3)))
		})
	})

})
