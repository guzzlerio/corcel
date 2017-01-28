package main

import (
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlan Assertions", func() {

	Context("ExactAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula").Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "talula"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", 1))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("EmptyAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "").Build()).
				WithAssertion(planBuilder.EmptyAssertion("value:1"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "1").Build()).
				WithAssertion(planBuilder.EmptyAssertion("value:1"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("GreaterThanAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.GreaterThanAssertion("value:1", 2))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
				WithAssertion(planBuilder.GreaterThanAssertion("value:1", 5))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("GreaterThanOrEqualAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.GreaterThanOrEqualAssertion("value:1", 5))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
				WithAssertion(planBuilder.GreaterThanOrEqualAssertion("value:1", 5))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("LessThanAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 3).Build()).
				WithAssertion(planBuilder.LessThanAssertion("value:1", 5))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.LessThanAssertion("value:1", 3))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("LessThanOrEqualAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.LessThanOrEqualAssertion("value:1", 5))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.LessThanOrEqualAssertion("value:1", 4))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("NotEmptyAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.NotEmptyAssertion("value:1"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:2", 5).Build()).
				WithAssertion(planBuilder.NotEmptyAssertion("value:1"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})

	Context("NotEqualAssertion", func() {

		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
				WithAssertion(planBuilder.NotEqualAssertion("value:1", 6))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})

		It("Fails", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", 6).Build()).
				WithAssertion(planBuilder.NotEqualAssertion("value:1", 6))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

	})
})
