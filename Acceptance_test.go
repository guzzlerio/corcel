package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	. "github.com/guzzlerio/corcel"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/rizo"

	. "github.com/guzzlerio/corcel/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance", func() {

	It("Outputs a summary to STDOUT", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
		}

		TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(500)
		})).For(rizo.RequestWithPath("/error"))

		output, err := SutExecute(list, "--summary")
		Expect(err).To(BeNil())

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Running Time: %v", summary.RunningTime)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Throughput: %.0f req/s", summary.Throughput)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Total Requests: %v", summary.TotalRequests)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Number of Errors: %v", summary.TotalErrors)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Availability: %v.0000%%", summary.Availability)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Bytes Sent: %v", summary.Bytes.TotalSent)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Bytes Received: %v", summary.Bytes.TotalReceived)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Mean Response Time: %.4f", summary.MeanResponseTime)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Min Response Time: %.4f ms", summary.MinResponseTime)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Max Response Time: %.4f ms", summary.MaxResponseTime)))
	})

	It("Error non-http url in the urls file causes a run time exception #21", func() {
		list := []string{
			fmt.Sprintf(`-Something`),
		}

		output, err := SutExecute(list)
		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring(errormanager.LogMessageVaidURLs))
	})

	It("Issue - Should write out panics to a log file and not std out", func() {
		planBuilder := yaml.NewPlanBuilder()

		planBuilder.
			SetIterations(1).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.IPanicAction().Build())

		output, err := ExecutePlanBuilder(planBuilder)
		Expect(err).ToNot(BeNil())

		Expect(string(output)).To(ContainSubstring("An unexpected error has occurred.  The error has been logged to /tmp/"))

		//Ensure that the file which was generated contains the error which caused the panic
		r, _ := regexp.Compile(`/tmp/[\w\d-]+`)
		var location = r.FindString(string(output))
		Expect(location).ToNot(Equal(""))
		data, err := ioutil.ReadFile(location)
		Expect(err).To(BeNil())
		Expect(string(data)).To(ContainSubstring("IPanicAction has caused this panic"))
	})
})
