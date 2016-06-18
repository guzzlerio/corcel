package main

import (
	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("ExecutionPlan Assertions", func() {

	Context("ExactAssertion", func() {

		It("Succeeds", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula").Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "talula"))

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
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", 1))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("EmptyAssertion", func() {

		It("Succeeds", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "").Build()).
				WithAssertion(planBuilder.EmptyAssertion("value:1"))

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
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "1").Build()).
				WithAssertion(planBuilder.EmptyAssertion("value:1"))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("GreaterThanAssertion", func() {

		It("Succeeds", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.GreaterThanAssertion("value:1", 2))

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
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
				WithAssertion(planBuilder.GreaterThanAssertion("value:1", 5))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("GreaterThanOrEqualAssertion", func() {

		It("Succeeds", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.GreaterThanOrEqualAssertion("value:1", 5))

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
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
				WithAssertion(planBuilder.GreaterThanOrEqualAssertion("value:1", 5))

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})
	/*

		Context("LessThanAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})

		Context("LessThanOrEqualAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})

		Context("NotEmptyAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})

		Context("NotEmptyAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})
	*/
})
