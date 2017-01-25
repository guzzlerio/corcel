package main_test

import (
	"fmt"
	"io/ioutil"
	"math"
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

	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.Header().Add("done", "1")
			w.WriteHeader(http.StatusOK)
		})

		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	Describe("Core Command Line Usage", func() {

		Describe("Iterations", func() {

			It("Plan Usage", func() {

				var plan = fmt.Sprintf(`---
name: Some Plan
iterations: 5
random: false
workers: 1
waitTime: 0s
duration: 0s
jobs:
    - name: Some Job
      steps:
      - name: Some Step
        action:
          type: HttpRequest
          httpHeaders:
            key: 1
          method: GET
          url: %s`, TestServer.CreateURL("/people"))

				summary, err := test.ExecutePlanFromData(plan, "--summary")

				Expect(err).To(BeNil())

				Expect(summary.TotalRequests).To(Equal(float64(5)))
			})

			It("List Usage", func() {
				list := []string{
					fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
				}
				summary, err := test.ExecuteList(list, "--summary", "--iterations", "5")
				Expect(err).To(BeNil())

				Expect(summary.TotalRequests).To(Equal(float64(5)))
			})
		})

		Describe("Workers", func() {

			It("Plan Usage", func() {

				var plan = fmt.Sprintf(`---
name: Some Plan
iterations: 0
random: false
workers: 5
waitTime: 0s
duration: 0s
jobs:
    - name: Some Job
      steps:
      - name: Some Step
        action:
          type: HttpRequest
          httpHeaders:
            key: 1
          method: GET
          url: %s`, TestServer.CreateURL("/people"))

				summary, err := test.ExecutePlanFromData(plan, "--summary")

				Expect(err).To(BeNil())

				Expect(summary.TotalRequests).To(Equal(float64(5)))
			})

			It("List Usage", func() {
				list := []string{
					fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
				}
				summary, err := test.ExecuteList(list, "--summary", "--workers", "5")
				Expect(err).To(BeNil())

				Expect(summary.TotalRequests).To(Equal(float64(5)))
			})
		})

		Describe("Wait Time", func() {

			It("Plan Usage", func() {

				var plan = fmt.Sprintf(`---
name: Some Plan
iterations: 5
random: false
workers: 1
waitTime: 1s
duration: 0s
jobs:
    - name: Some Job
      steps:
      - name: Some Step
        action:
          type: HttpRequest
          httpHeaders:
            key: 1
          method: GET
          url: %s`, TestServer.CreateURL("/people"))

				summary, err := test.ExecutePlanFromData(plan, "--summary")

				Expect(err).To(BeNil())

				Expect(math.Floor(summary.RunningTime.Seconds())).To(Equal(float64(5)))
			})

			It("List Usage", func() {
				list := []string{
					fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
				}
				summary, err := test.ExecuteList(list, "--summary", "--wait-time", "1s", "--iterations", "5")
				Expect(err).To(BeNil())

				Expect(math.Floor(summary.RunningTime.Seconds())).To(Equal(float64(5)))
			})
		})

		Describe("Duration", func() {

			It("Plan Usage", func() {

				var plan = fmt.Sprintf(`---
name: Some Plan
iterations: 0
random: false
workers: 1
waitTime: 0s
duration: 5s
jobs:
    - name: Some Job
      steps:
      - name: Some Step
        action:
          type: HttpRequest
          httpHeaders:
            key: 1
          method: GET
          url: %s`, TestServer.CreateURL("/people"))

				summary, err := test.ExecutePlanFromData(plan, "--summary")

				Expect(err).To(BeNil())

				Expect(math.Floor(summary.RunningTime.Seconds())).To(Equal(float64(5)))
			})

			It("List Usage", func() {
				list := []string{
					fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
				}
				summary, err := test.ExecuteList(list, "--summary", "--duration", "5s")
				Expect(err).To(BeNil())

				Expect(math.Floor(summary.RunningTime.Seconds())).To(Equal(float64(5)))
			})
		})
	})

	It("Halts execution if a payload input file is not found", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST -d '{"name":"talula"}'`, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X POST -d @missing-file.json`, URLForTestServer("/success")),
		}

		summary, err := test.ExecuteList(list, "--summary")
		Expect(err).ToNot(BeNil())

		Expect(summary.Error).To(ContainSubstring("Request body file not found: missing-file.json"))
	})

	It("Error non-http url in the urls file causes a run time exception #21", func() {
		list := []string{
			fmt.Sprintf(`-Something`),
		}

		summary, err := test.ExecuteList(list, "--summary")
		Expect(err).ToNot(BeNil())
		Expect(summary.Error).To(ContainSubstring(errormanager.LogMessageVaidURLs))
	})

	It("Issue - Should write out panics to a log file and not std out", func() {
		planBuilder := yaml.NewPlanBuilder()

		planBuilder.
			SetIterations(1).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.IPanicAction().Build())

		output, err := test.ExecutePlanBuilder(planBuilder)
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
