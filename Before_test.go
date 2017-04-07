package main

import (
	"net/http"
	"testing"

	"github.com/guzzlerio/rizo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
)

func TestBefore_After(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("Before After", t, func() {
		var (
			planBuilder *yaml.PlanBuilder
			path        string
			body        string
		)

		func() {
			TestServer.Clear()
			factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
			})

			TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
			planBuilder = yaml.NewPlanBuilder()
			path = "/people"
			body = "Zee Body"
		}()

		getBody := func(requests []rizo.RecordedRequest) []string {
			var bodies []string
			for _, request := range requests {
				bodies = append(bodies, request.Body)
			}
			return bodies
		}

		Convey("Plan", func() {
			Convey("Before hook", func() {
				Convey("is invoked before plan execution", func() {
					planBuilder.
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan").Build()).
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Before Plan",
						"Zee Body",
					})
				})

				Convey("with multiple actions is invoked in order before plan execution", func() {
					planBuilder.
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan 1").Build()).
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan 2").Build()).
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Before Plan 1",
						"Before Plan 2",
						"Zee Body",
					})
				})
			})

			Convey("after hook is invoked after plan execution", func() {
				planBuilder.
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Plan").Build()).
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_, _ = test.ExecutePlanBuilder(planBuilder)
				So(getBody(TestServer.Requests), ShouldResemble, []string{
					"Zee Body",
					"After Plan",
				})
			})

			Convey("before and after hooks are invoked before and after plan execution", func() {
				planBuilder.
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan").Build()).
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Plan").Build()).
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_, _ = test.ExecutePlanBuilder(planBuilder)
				So(getBody(TestServer.Requests), ShouldResemble, []string{
					"Before Plan",
					"Zee Body",
					"After Plan",
				})
			})
		})

		Convey("Job", func() {
			Convey("Before hook", func() {
				Convey("is invoked before job execution", func() {
					planBuilder.
						CreateJob().
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job").Build()).
						CreateStep().
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Before Job",
						"Zee Body",
					})
				})
			})

			Convey("After hook", func() {
				Convey("is invoked after job execution", func() {
					planBuilder.
						CreateJob().
						AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job").Build()).
						CreateStep().
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Zee Body",
						"After Job",
					})
				})
			})

			Convey("Before and After hook", func() {
				Convey("is invoked before and after job execution", func() {
					planBuilder.
						CreateJob().
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job").Build()).
						AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job").Build()).
						CreateStep().
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Before Job",
						"Zee Body",
						"After Job",
					})
				})
			})
		})

		Convey("Step", func() {
			Convey("Before hook", func() {
				Convey("is invoked before step execution", func() {
					jobBuilder := planBuilder.CreateJob()

					jobBuilder.CreateStep().
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					jobBuilder.CreateStep().
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Step").Build()).
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Zee Body",
						"Before Step",
						"Zee Body",
					})
				})
			})

			Convey("After hook", func() {
				Convey("is invoked after job execution", func() {
					jobBuilder := planBuilder.CreateJob()

					jobBuilder.CreateStep().
						AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 1").Build()).
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					jobBuilder.CreateStep().
						AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 2").Build()).
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Zee Body",
						"After Step 1",
						"Zee Body",
						"After Step 2",
					})
				})
			})

			Convey("Before and After hook", func() {
				Convey("is invoked before and after job execution", func() {
					jobBuilder := planBuilder.CreateJob()

					jobBuilder.CreateStep().
						AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 1").Build()).
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					jobBuilder.CreateStep().
						AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Step 2").Build()).
						AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 2").Build()).
						ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

					_, _ = test.ExecutePlanBuilder(planBuilder)
					So(getBody(TestServer.Requests), ShouldResemble, []string{
						"Zee Body",
						"After Step 1",
						"Before Step 2",
						"Zee Body",
						"After Step 2",
					})
				})
			})
		})

		Convey("Combined", func() {
			Convey("handles a mix", func() {
				planBuilder.
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan").Build())

				jobBuilder := planBuilder.CreateJob("Job 1").
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job 1").Build())

				jobBuilder.CreateStep().
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job 1, Step 1").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				jobBuilder.CreateStep().
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job 1, Step 2").Build()).
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job 1, Step 2").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				jobBuilder = planBuilder.CreateJob("Job 2").
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job 2").Build())

				jobBuilder.CreateStep().
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job 2, Step 1").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())
				_, _ = test.ExecutePlanBuilder(planBuilder)
				So(getBody(TestServer.Requests), ShouldResemble, []string{
					"Before Plan",
					"Before Job 1",
					"Zee Body",
					"After Job 1, Step 1",
					"Before Job 1, Step 2",
					"Zee Body",
					"After Job 1, Step 2",
					"Before Job 2",
					"Zee Body",
					"After Job 2, Step 1",
				})
			})
		})
	})
}
