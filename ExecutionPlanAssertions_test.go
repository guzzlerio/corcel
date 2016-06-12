package main

import (
	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = PDescribe("ExecutionPlan Assertions", func() {

	Context("ExactAssertion", func() {

		It("Succeeds", func() {

			planBuilder := test.NewYamlPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(map[string]interface{}{
				"type": "DummyAction",
				"results": map[string]string{
					"value:1": "talula",
				},
			})

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
		})

		It("Fails", func() {

		})

	})
	/*

		Context("EmptyAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})

		Context("GreaterThanAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})

		Context("GreaterThanOrEqualAssertion", func() {

			It("Succeeds", func() {

			})

			It("Fails", func() {

			})

		})

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
