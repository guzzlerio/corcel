package main_test

import (
	. "ci.guzzler.io/guzzler/corcel"
	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanContexts", func() {

	It("Plan Scope", func() {
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
})
