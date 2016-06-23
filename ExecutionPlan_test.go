package main

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/global"
	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"
)

var _ = Describe("ExecutionPlan", func() {

	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	/*/\/\/\/\/\/\/\//\/\/\/\/\/\/\/\/\/\/\/\\/\/
	 *
	 * REFER TO ISSUE #50
	 *
	 *\/\/\/\/\/\/\//\/\/\/\/\/\/\/\/\/\/\/\\/\/
	 */

	for _, numberOfWorkers := range global.NumberOfWorkersToTest {
		name := fmt.Sprintf("SetWorkers for %v workers", numberOfWorkers)
		func(workers int) {
			It(name, func() {
				planBuilder := test.NewYamlPlanBuilder()

				planBuilder.
					SetWorkers(workers).
					CreateJob().
					CreateStep().
					ToExecuteAction(GetHTTPRequestAction("/people"))

				err := ExecutePlanBuilder(planBuilder)
				Expect(err).To(BeNil())

				var executionOutput statistics.AggregatorSnapShot
				utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
				var summary = statistics.CreateSummary(executionOutput)

				Expect(summary.TotalErrors).To(Equal(float64(0)))
				Expect(summary.TotalRequests).To(Equal(float64(workers)))
				Expect(len(TestServer.Requests)).To(Equal(workers))
			})
		}(numberOfWorkers)
	}

	It("SetWaitTime", func() {
		numberOfSteps := 6
		waitTime := "500ms"

		planBuilder := test.NewYamlPlanBuilder()
		planBuilder.SetWaitTime(waitTime)
		jobBuilder := planBuilder.CreateJob()

		for i := 0; i < numberOfSteps; i++ {
			jobBuilder.CreateStep().ToExecuteAction(GetHTTPRequestAction("/people"))
		}

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		var executionOutput statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		actual, _ := time.ParseDuration(summary.RunningTime)
		seconds := actual.Seconds()
		seconds = math.Floor(seconds)
		Expect(seconds).To(Equal(float64(3)))
	})

	It("SetDuration", func() {
		duration := "3s"
		planBuilder := test.NewYamlPlanBuilder()
		planBuilder.SetDuration(duration)
		jobBuilder := planBuilder.CreateJob()
		jobBuilder.CreateStep().ToExecuteAction(GetHTTPRequestAction("/people"))
		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		var executionOutput statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		actual, _ := time.ParseDuration(summary.RunningTime)
		seconds := actual.Seconds()
		seconds = math.Floor(seconds)
		Expect(seconds).To(Equal(float64(3)))
	})

	It("SetRandom", func() {
		numberOfSteps := 6

		planBuilder := test.NewYamlPlanBuilder()
		planBuilder.SetRandom(true)

		for i := 0; i < numberOfSteps; i++ {
			jobBuilder := planBuilder.CreateJob()
			jobBuilder.CreateStep().ToExecuteAction(GetHTTPRequestAction(fmt.Sprintf("/%d", i+1)))
		}

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())
		firstBatchOfRequests := utils.ConcatRequestPaths(utils.ToHTTPRequestArray(TestServer.Requests))
		TestServer.Clear()

		err = ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())
		secondBatchOfRequests := utils.ConcatRequestPaths(utils.ToHTTPRequestArray(TestServer.Requests))

		Expect(firstBatchOfRequests).ToNot(Equal(secondBatchOfRequests))
	})

	PDescribe("HttpRequest", func() {})

	Describe("Assertions", func() {

		BeforeEach(func() {
			TestServer.Clear()
			TestServer.Use(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
			}).For(rizo.RequestWithPath("/boom"))
		})
		//ASSERTION FAILURES ARE NOT CURRENTLY COUNTING AS ERRORS IN THE SUMMARY OUTPUT
		It("ExactAssertion Fails", func() {

			planBuilder := test.NewYamlPlanBuilder()
			planBuilder.CreateJob().
				CreateStep().
				ToExecuteAction(GetHTTPRequestAction("/boom")).
				WithAssertion(HTTPStatusExactAssertion(201))

			err := ExecutePlanBuilder(planBuilder)
			utils.CheckErr(err)

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertions).To(Equal(int64(1)))
			Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
		})

		It("ExactAssertion Pass", func() {

			planBuilder := test.NewYamlPlanBuilder()
			planBuilder.CreateJob().
				CreateStep().
				ToExecuteAction(GetHTTPRequestAction("/boom")).
				WithAssertion(HTTPStatusExactAssertion(200))

			err := ExecutePlanBuilder(planBuilder)
			utils.CheckErr(err)

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalAssertions).To(Equal(int64(1)))
			Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
		})
	})

	It("Name", func() {
		planBuilder := test.NewYamlPlanBuilder()
		planBuilder.CreateJob().
			CreateStep().
			ToExecuteAction(GetHTTPRequestAction("/people")).
			WithAssertion(HTTPStatusExactAssertion(201))

		err := ExecutePlanBuilder(planBuilder)
		utils.CheckErr(err)

		var executionOutput statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(summary.TotalRequests).To(BeNumerically(">", 0))
		Expect(summary.TotalErrors).To(Equal(float64(0)))
	})
})
