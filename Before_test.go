package main

import (
	"net/http"

	"ci.guzzler.io/guzzler/corcel/test"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Before After", func() {
	var (
		planBuilder *test.YamlPlanBuilder
		path        string
		body        string
	)
	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
		planBuilder = test.NewYamlPlanBuilder()
		path = "/people"
		body = "Zee Body"
	})

	getBody := func(requests []rizo.RecordedRequest) []string {
		var bodies []string
		for _, request := range requests {
			bodies = append(bodies, request.Body)
		}
		return bodies
	}

	Context("Plan", func() {
		Context("Before hook", func() {
			It("is invoked before plan execution", func() {
				planBuilder.
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan").Build()).
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Before Plan",
					"Zee Body",
				})))
			})

			It("with multiple actions is invoked in order before plan execution", func() {
				planBuilder.
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan 1").Build()).
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan 2").Build()).
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Before Plan 1",
					"Before Plan 2",
					"Zee Body",
				})))
			})
		})

		It("after hook is invoked after plan execution", func() {
			planBuilder.
				AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Plan").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

			_ = test.ExecutePlanBuilder("./corcel", planBuilder)
			Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
				"Zee Body",
				"After Plan",
			})))
		})

		It("before and after hooks are invoked before and after plan execution", func() {
			planBuilder.
				AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Plan").Build()).
				AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Plan").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

			_ = test.ExecutePlanBuilder("./corcel", planBuilder)
			Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
				"Before Plan",
				"Zee Body",
				"After Plan",
			})))
		})
	})

	Context("Job", func() {
		Context("Before hook", func() {
			It("is invoked before job execution", func() {
				planBuilder.
					CreateJob().
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job").Build()).
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Before Job",
					"Zee Body",
				})))
			})
		})

		Context("After hook", func() {
			It("is invoked after job execution", func() {
				planBuilder.
					CreateJob().
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job").Build()).
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Zee Body",
					"After Job",
				})))
			})
		})

		Context("Before and After hook", func() {
			It("is invoked before and after job execution", func() {
				planBuilder.
					CreateJob().
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Job").Build()).
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Job").Build()).
					CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Before Job",
					"Zee Body",
					"After Job",
				})))
			})
		})
	})

	Context("Step", func() {
		Context("Before hook", func() {
			It("is invoked before step execution", func() {
				jobBuilder := planBuilder.CreateJob()

				jobBuilder.CreateStep().
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				jobBuilder.CreateStep().
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Step").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Zee Body",
					"Before Step",
					"Zee Body",
				})))
			})
		})

		Context("After hook", func() {
			It("is invoked after job execution", func() {
				jobBuilder := planBuilder.CreateJob()

				jobBuilder.CreateStep().
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 1").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				jobBuilder.CreateStep().
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 2").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Zee Body",
					"After Step 1",
					"Zee Body",
					"After Step 2",
				})))
			})
		})

		Context("Before and After hook", func() {
			It("is invoked before and after job execution", func() {
				jobBuilder := planBuilder.CreateJob()

				jobBuilder.CreateStep().
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 1").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				jobBuilder.CreateStep().
					AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("Before Step 2").Build()).
					AddAfter(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body("After Step 2").Build()).
					ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

				_ = test.ExecutePlanBuilder("./corcel", planBuilder)
				Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
					"Zee Body",
					"After Step 1",
					"Before Step 2",
					"Zee Body",
					"After Step 2",
				})))
			})
		})
	})

	Context("Combined", func() {
		It("handles a mix", func() {
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
			_ = test.ExecutePlanBuilder("./corcel", planBuilder)
			Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{
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
			})))
		})
	})
})
