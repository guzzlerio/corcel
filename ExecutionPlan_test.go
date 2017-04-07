package main

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/guzzlerio/rizo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	"github.com/guzzlerio/corcel/utils"
)

func TestExecutionPlan(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("ExecutionPlan", t, func() {

		func() {
			TestServer.Clear()
			factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
			})

			TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
		}()

		defer func() {
			TestServer.Clear()
		}()

		Convey("SetIterations", func() {
			Convey("Single Job Single Step", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					SetIterations(2).
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Build())

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalRequests, ShouldEqual, float64(2))
			})
			Convey("Single Job Two Steps", func() {
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

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalRequests, ShouldEqual, float64(4))
			})
			Convey("Two Jobs Two Steps", func() {
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

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalRequests, ShouldEqual, float64(8))
			})

			Convey("For a list", func() {
				list := []string{
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				}

				summary, err := test.ExecuteListForApplication(list, config.Configuration{
					Iterations: 5,
				})
				So(err, ShouldBeNil)
				So(summary.TotalRequests, ShouldEqual, float64(30))
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
				Convey(name, func() {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						SetWorkers(workers).
						CreateJob().
						CreateStep().
						ToExecuteAction(GetHTTPRequestAction("/people"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					So(err, ShouldBeNil)
					So(summary.TotalErrors, ShouldEqual, float64(0))
					So(summary.TotalRequests, ShouldEqual, float64(workers))
					So(len(TestServer.Requests), ShouldEqual, workers)
				})
			}(numberOfWorkers)
		}

		Convey("SetWaitTime", func() {
			numberOfSteps := 6
			waitTime := "500ms"

			planBuilder := yaml.NewPlanBuilder()
			planBuilder.SetWaitTime(waitTime)
			jobBuilder := planBuilder.CreateJob()

			for i := 0; i < numberOfSteps; i++ {
				jobBuilder.CreateStep().ToExecuteAction(GetHTTPRequestAction("/people"))
			}

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			So(err, ShouldBeNil)
			actual := summary.RunningTime
			seconds := actual.Seconds()
			seconds = math.Floor(seconds)
			So(seconds, ShouldEqual, float64(3))
		})

		Convey("SetDuration", func() {
			duration := "3s"
			planBuilder := yaml.NewPlanBuilder()
			planBuilder.SetDuration(duration)
			jobBuilder := planBuilder.CreateJob()
			jobBuilder.CreateStep().ToExecuteAction(GetHTTPRequestAction("/people"))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			So(err, ShouldBeNil)
			actual := summary.RunningTime
			seconds := actual.Seconds()
			seconds = math.Floor(seconds)
			So(seconds, ShouldEqual, float64(3))
		})

		Convey("SetRandom", func() {
			numberOfSteps := 6

			planBuilder := yaml.NewPlanBuilder()
			planBuilder.SetRandom(true)

			for i := 0; i < numberOfSteps; i++ {
				jobBuilder := planBuilder.CreateJob()
				jobBuilder.CreateStep().ToExecuteAction(GetHTTPRequestAction(fmt.Sprintf("/%d", i+1)))
			}

			_, err := test.ExecutePlanBuilderForApplication(planBuilder)
			So(err, ShouldBeNil)

			firstBatchOfRequests := utils.ConcatRequestPaths(utils.ToHTTPRequestArray(TestServer.Requests))
			TestServer.Clear()

			_, err = test.ExecutePlanBuilderForApplication(planBuilder)
			So(err, ShouldBeNil)
			secondBatchOfRequests := utils.ConcatRequestPaths(utils.ToHTTPRequestArray(TestServer.Requests))

			So(firstBatchOfRequests, ShouldNotEqual, secondBatchOfRequests)
		})

		SkipConvey("HttpRequest", func() {})

		Convey("Assertions", func() {

			func() {
				TestServer.Clear()
				TestServer.Use(func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusOK)
				}).For(rizo.RequestWithPath("/boom"))
			}()

			//ASSERTION FAILURES ARE NOT CURRENTLY COUNTING AS ERRORS IN THE SUMMARY OUTPUT
			Convey("ExactAssertion Fails", func() {

				planBuilder := yaml.NewPlanBuilder()
				planBuilder.CreateJob().
					CreateStep().
					ToExecuteAction(GetHTTPRequestAction("/boom")).
					WithAssertion(HTTPStatusExactAssertion(201))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertions, ShouldEqual, int64(1))
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

			Convey("ExactAssertion Pass", func() {

				planBuilder := yaml.NewPlanBuilder()
				planBuilder.CreateJob().
					CreateStep().
					ToExecuteAction(GetHTTPRequestAction("/boom")).
					WithAssertion(HTTPStatusExactAssertion(200))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertions, ShouldEqual, int64(1))
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})
		})

		Convey("Name", func() {
			planBuilder := yaml.NewPlanBuilder()
			planBuilder.CreateJob().
				CreateStep().
				ToExecuteAction(GetHTTPRequestAction("/people")).
				WithAssertion(HTTPStatusExactAssertion(201))

			summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
			So(err, ShouldBeNil)
			So(summary.TotalRequests, ShouldBeGreaterThan, 0)
			So(summary.TotalErrors, ShouldEqual, float64(0))
		})
	})
}
