package main

import (
	"net/http"

	"ci.guzzler.io/guzzler/corcel/test"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Before After", func() {
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
		PContext("After hook", func() {
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
	})
})
