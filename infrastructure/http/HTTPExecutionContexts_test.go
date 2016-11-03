package http_test

import (
	"net/http"
	"strings"

	"github.com/guzzlerio/corcel/test"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanContexts", func() {

	path := "/something"
	BeforeEach(func() {
		TestServer.Clear()

		TestServer.Use(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		}).For(rizo.RequestWithPath(path))
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	Context("Using List Variables", func() {
		It("inside the http headers", func() {

			expectedHeaderKey := "Content-Type"
			json := "application/json"
			xml := "application/json"
			carf := "application/carf"

			planBuilder := test.NewYamlPlanBuilder()
			planBuilder.
				SetIterations(3).
				WithContext(planBuilder.BuildContext().SetList("Content-type", []map[string]interface{}{
				{"commonType": json},
				{"commonType": xml},
				{"commonType": carf},
			}).Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().Header(expectedHeaderKey, "$Content-type.commonType").URL(TestServer.CreateURL(path)).Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			Expect(len(TestServer.Requests)).To(Equal(3))
			Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(expectedHeaderKey, json))).To(Equal(true))
			Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(expectedHeaderKey, xml))).To(Equal(true))
			Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(expectedHeaderKey, carf))).To(Equal(true))

		})
	})

	Context("Using variables", func() {

		It("inside the http headers", func() {

			expectedHeaderKey := "Content-Type"
			expectedHeaderValue := "application/json"

			planBuilder := test.NewYamlPlanBuilder()
			planBuilder.WithContext(planBuilder.BuildContext().Set("commonType", expectedHeaderValue).Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().Header(expectedHeaderKey, "$commonType").URL(TestServer.CreateURL(path)).Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(expectedHeaderKey, expectedHeaderValue))).To(Equal(true))
		})

		It("inside the url", func() {

			path := "/$path?a=$a&b=$b&c=$c"

			planBuilder := test.NewYamlPlanBuilder()
			planBuilder.WithContext(planBuilder.BuildContext().Set("path", "fubar").Set("a", "1").Set("b", "2").Set("c", "3").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			Expect(TestServer.Find(rizo.RequestWithPath("/fubar"), rizo.RequestWithQuerystring("a=1&b=2&c=3"))).To(Equal(true))
		})

		It("inside the body", func() {
			body := `
        {
          "firstname" : "$firstname",
          "lastname" : "$lastname"
        }
      `
			planBuilder := test.NewYamlPlanBuilder()
			planBuilder.WithContext(planBuilder.BuildContext().Set("firstname", "john").Set("lastname", "doe").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().Header("Content-type", "application/json").Body(body).URL(TestServer.CreateURL(path)).Build())

			err := ExecutePlanBuilder(planBuilder)
			Expect(err).To(BeNil())

			expectedBody := strings.Replace(body, "$firstname", "john", -1)
			expectedBody = strings.Replace(expectedBody, "$lastname", "doe", -1)

			Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithBody(expectedBody))).To(Equal(true))
		})

	})

	It("Set the QueryString", func() {

	})

	It("Extend the QueryString", func() {
		//If a base querystring is set the jobs, steps and actions add/override the previous
	})

	It("Set HTTP Header", func() {

		planBuilder := test.NewYamlPlanBuilder()

		expectedHeaderKey := "content-boomboom"
		expectedHeaderValue := "bang/boom"
		headers := map[string]string{}
		headers[expectedHeaderKey] = expectedHeaderValue

		planBuilder.WithContext(planBuilder.BuildContext().Set("httpHeaders", headers).Build()).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Build())

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(expectedHeaderKey, expectedHeaderValue))).To(Equal(true))
	})

	It("Context does not override a HTTP Header set in the action it self", func() {
		planBuilder := test.NewYamlPlanBuilder()

		contextHeaderKey := "content-boomboom"
		contextHeaderValue := "bang/boom"
		headers := map[string]string{}
		headers[contextHeaderKey] = contextHeaderValue

		expectedHeaderValue := "hazaa"

		planBuilder.WithContext(planBuilder.BuildContext().Set("httpHeaders", headers).Build()).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Header(contextHeaderKey, expectedHeaderValue).Build())

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(contextHeaderKey, expectedHeaderValue))).To(Equal(true))
	})
})
