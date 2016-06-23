package main_test

import (
	"net/http"

	. "ci.guzzler.io/guzzler/corcel"
	"ci.guzzler.io/guzzler/corcel/test"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("ExecutionPlanHttpRequest", func() {
	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
	})

	It("Supplies a payload to the HTTP Request", func() {
		planBuilder := test.NewYamlPlanBuilder()

		path := "/people"
		body := "Zee Body"

		planBuilder.
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.HTTPRequestAction().URL(TestServer.CreateURL(path)).Body(body).Build())

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())
		Expect(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithBody(body))).To(Equal(true))
	})
})
