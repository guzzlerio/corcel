package main

import (
	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanContexts", func() {

	Context("Plan Scope", func() {
		It("Succeeds", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("value:3", "3"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})
		It("Fails", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("value:3", "3"))

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
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build()).
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("value:3", "3"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build())

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("value:3", "3"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(3)))
		})
	})

})
