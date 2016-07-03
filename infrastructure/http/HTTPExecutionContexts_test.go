package http_test

import (
	"net/http"

	"ci.guzzler.io/guzzler/corcel/test"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanContexts", func() {

	BeforeEach(func() {
		TestServer.Clear()

	})

	AfterEach(func() {
		TestServer.Clear()
	})

	It("Set HTTP Header", func() {
		path := "/something"

		TestServer.Use(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		}).For(rizo.RequestWithPath(path))

		planBuilder := test.NewYamlPlanBuilder()

		expectedHeaderKey := "content-boomboom"
		expectedHeaderValue := "bang/boom"
		headers := map[string]string{}
		headers[expectedHeaderKey] = expectedHeaderValue

		planBuilder.WithContext(planBuilder.BuildContext().Set("httpHeaders", headers).Build()).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPRequestAction().URL(TestServer.CreateURL(path)).Build())

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(expectedHeaderKey, expectedHeaderValue))).To(Equal(true))
	})

	It("Context does not override a HTTP Header set in the action it self", func() {
		path := "/something"

		TestServer.Use(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		}).For(rizo.RequestWithPath(path))

		planBuilder := test.NewYamlPlanBuilder()

		contextHeaderKey := "content-boomboom"
		contextHeaderValue := "bang/boom"
		headers := map[string]string{}
		headers[contextHeaderKey] = contextHeaderValue

		expectedHeaderValue := "hazaa"

		planBuilder.WithContext(planBuilder.BuildContext().Set("httpHeaders", headers).Build()).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPRequestAction().URL(TestServer.CreateURL(path)).Header(contextHeaderKey, expectedHeaderValue).Build())

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithHeader(contextHeaderKey, expectedHeaderValue))).To(Equal(true))
	})
})
