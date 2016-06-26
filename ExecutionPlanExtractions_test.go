package main_test

import (
	. "ci.guzzler.io/guzzler/corcel"
	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanExtractions", func() {
	Context("Regex", func() {
		Context("Step Scope", func() {
			It("Matches simple pattern", func() {
				planBuilder := test.NewYamlPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
					WithExtractor(planBuilder.RegexExtractor().Name("regex:match:1").Key("value:1").Match("\\d+").Build()).
					WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

				err := ExecutePlanBuilder(planBuilder)
				Expect(err).To(BeNil())

				var executionOutput statistics.AggregatorSnapShot
				utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
				var summary = statistics.CreateSummary(executionOutput)

				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			})

			PIt("Extends the name with any named groups", func() {})

			PIt("Extends the name with index access with any non-named groups", func() {})
		})
		PContext("Job Scope", func() {

		})
		PContext("Plan Scope", func() {

		})
	})

	Context("XPAth", func() {
		PContext("Step Scope", func() {

		})
		PContext("Job Scope", func() {

		})
		PContext("Plan Scope", func() {

		})
	})

	Context("JSON Path", func() {
		PContext("Step Scope", func() {

		})
		PContext("Job Scope", func() {

		})
		PContext("Plan Scope", func() {

		})
	})

	Context("Javascript", func() {
		PContext("Step Scope", func() {

		})
		PContext("Job Scope", func() {

		})
		PContext("Plan Scope", func() {

		})
	})
})
