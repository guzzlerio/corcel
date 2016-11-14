package main

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/utils"
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

	Context("SetIterations", func() {
		It("Single Job Single Step", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				SetIterations(2).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalRequests).To(Equal(float64(2)))
		})
		It("Single Job Two Steps", func() {
			planBuilder := yaml.NewPlanBuilder()

			jobBuilder := planBuilder.
				SetIterations(2).
				CreateJob()
			jobBuilder.
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())
			jobBuilder.
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalRequests).To(Equal(float64(4)))
		})
		It("Single Job Two Steps", func() {
			planBuilder := yaml.NewPlanBuilder()

			jobBuilder := planBuilder.
				SetIterations(2).
				CreateJob()

			jobBuilder.
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())
			jobBuilder.
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())

			jobBuilder = planBuilder.
				CreateJob()
			jobBuilder.
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())
			jobBuilder.
				CreateStep().
				ToExecuteAction(planBuilder.DummyAction().Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalRequests).To(Equal(float64(8)))
		})

		It("For a list", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			}

			SutExecute(list, "--iterations", "5")

			var executionOutput statistics.AggregatorSnapShot
			utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalRequests).To(Equal(float64(30)))
		})

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
				planBuilder := yaml.NewPlanBuilder()

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

		planBuilder := yaml.NewPlanBuilder()
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
		planBuilder := yaml.NewPlanBuilder()
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

		planBuilder := yaml.NewPlanBuilder()
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

			planBuilder := yaml.NewPlanBuilder()
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

			planBuilder := yaml.NewPlanBuilder()
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
		planBuilder := yaml.NewPlanBuilder()
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
