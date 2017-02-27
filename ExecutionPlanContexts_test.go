package main

import (
	"os"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
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

				_, err := test.ExecutePlanBuilderForApplication(planBuilder)
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

			jobBuilder := planBuilder.
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build()).
				CreateJob()

			jobBuilder.CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build())

			jobBuilder.CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build()).
				WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("$value:3", "3")).
				WithAssertion(planBuilder.ExactAssertion("something", "1"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
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

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(3)))
		})
	})

	Context("Job Scope", func() {
		It("Succeeds", func() {
			planBuilder := yaml.NewPlanBuilder()

			jobBuilder := planBuilder.
				CreateJob().
				WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build())

			jobBuilder.CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build())

			jobBuilder.CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build()).
				WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
				WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
				WithAssertion(planBuilder.ExactAssertion("$value:3", "3")).
				WithAssertion(planBuilder.ExactAssertion("something", "1"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
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

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			Expect(err).To(BeNil())
			Expect(summary.TotalAssertionFailures).To(Equal(int64(3)))
		})
	})

})
