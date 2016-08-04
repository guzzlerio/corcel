package main

import (
	// "ci.guzzler.io/guzzler/corcel/statistics"

	"fmt"
	"net/http"

	"ci.guzzler.io/guzzler/corcel/test"
	// "ci.guzzler.io/guzzler/corcel/utils"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Before", func() {
	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
	})

	Context("Plan", func() {
		It("before hook is invoked before plan execution", func() {
			planBuilder := test.NewYamlPlanBuilder()
			path := "/people"
			body := "Zee Body"

			planBuilder.
				AddBefore(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Header("before", "plan").Body("Before Plan").Build()).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

			err := test.ExecutePlanBuilder("./corcel", planBuilder)
			Expect(err).To(BeNil())
			fmt.Printf("%+v", TestServer.Requests)
			Expect(TestServer.Requests).Should(HaveLen(2))
			getBody := func(requests []rizo.RecordedRequest) []string {
				var bodies []string
				for _, request := range requests {
					bodies = append(bodies, request.Body)
				}
				return bodies
			}
			Expect(TestServer.Requests).Should(WithTransform(getBody, Equal([]string{"Before Plan", "Zee Body"})))
		})
	})
})
