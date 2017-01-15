package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	. "github.com/guzzlerio/corcel"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	"github.com/guzzlerio/rizo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance", func() {

	It("Halts execution if a payload input file is not found", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST -d '{"name":"talula"}'`, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X POST -d @missing-file.json`, URLForTestServer("/success")),
		}

		output, err := test.ExecuteList("./corcel", list)
		Expect(err).ToNot(BeNil())

		Expect(string(output)).To(ContainSubstring("Request body file not found: missing-file.json"))
	})

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

		output, err := test.ExecuteList("./corcel", list, "--summary")
		Expect(err).To(BeNil())

		Expect(string(output)).To(ContainSubstring("Summary"))
	})

	It("Error non-http url in the urls file causes a run time exception #21", func() {
		list := []string{
			fmt.Sprintf(`-Something`),
		}

		output, err := test.ExecuteList("./corcel", list)
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

	Describe("Shell", func() {

		Describe("run", func() {

			PIt("setting iterations", func() {

			})

			PIt("setting duration", func() {

			})

			PIt("setting workers", func() {

			})

			PIt("setting plan", func() {

			})

			PIt("setting summary", func() {

			})
		})

	})
})
