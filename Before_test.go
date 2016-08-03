package main

import (
	// "ci.guzzler.io/guzzler/corcel/statistics"

	"ci.guzzler.io/guzzler/corcel/test"
	// "ci.guzzler.io/guzzler/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Before", func() {
	Context("Plan", func() {
		FIt("before hook is invoked before plan execution", func() {
			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				AddBefore(planBuilder.DummyAction().Set("before:plan", "executed before").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula").Build()).
				WithAssertion(planBuilder.ExactAssertion("value:1", "talula"))

			err := test.ExecutePlanBuilder("./corcel", planBuilder)
			Expect(err).To(BeNil())

			// var executionOutput statistics.AggregatorSnapShot
			// utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			// var summary = statistics.CreateSummary(executionOutput)

			// So(summary.TotalAssertionFailures, ShouldEqual, 0)
		})
	})
})
