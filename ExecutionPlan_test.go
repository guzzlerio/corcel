package main

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/REAANDREW/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/processor"
	"ci.guzzler.io/guzzler/corcel/statistics"
	. "ci.guzzler.io/guzzler/corcel/utils"
)

var _ = Describe("ExecutionPlan", func() {

	/*/\/\/\/\/\/\/\//\/\/\/\/\/\/\/\/\/\/\/\\/\/
	 *
	 * REFER TO ISSUE #50
	 *
	 *\/\/\/\/\/\/\//\/\/\/\/\/\/\/\/\/\/\/\\/\/
	 */

	It("Name", func() {
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))

		planBuilder := processor.NewYamlPlanBuilder()
		planBuilder.CreateJob().CreateStep().Action(map[string]interface{}{
			"type":          "HttpRequest",
			"requesTimeout": 150,
			"method":        "GET",
			"url":           TestServer.CreateURL("/people"),
			"httpHeaders": map[string]string{
				"Content-Type": "application/json",
			},
		}).Assertion(map[string]interface{}{
			"type":     "ExactAssertion",
			"key":      "http:response:status",
			"expected": 202,
		})
		file, err := planBuilder.Build()
		CheckErr(err)
		exePath, err := filepath.Abs("./corcel")
		CheckErr(err)
		defer func() {
			err := os.Remove(file.Name())
			CheckErr(err)
		}()
		args := []string{"--plan"}
		cmd := exec.Command(exePath, append(args, file.Name())...)
		_, err = cmd.CombinedOutput()
		CheckErr(err)
		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)
		Expect(summary.TotalRequests).To(BeNumerically(">", 0))
		Expect(summary.TotalErrors).To(Equal(float64(0)))
	})
})
