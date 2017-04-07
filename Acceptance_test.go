package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"testing"

	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	"github.com/guzzlerio/rizo"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAcceptance(t *testing.T) {

	BeforeTest()

	defer AfterTest()

	Convey("Acceptance", t, func() {

		func() {
			TestServer.Clear()
			factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				w.Header().Add("done", "1")
				w.WriteHeader(http.StatusOK)
			})

			TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
		}()

		defer func() {
			TestServer.Clear()
		}()

		Convey("Core Command Line Usage", func() {

			Convey("Iterations", func() {

				Convey("Plan Usage", func() {

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

					So(err, ShouldBeNil)

					So(summary.TotalRequests, ShouldEqual, float64(5))
				})

				Convey("List Usage", func() {
					list := []string{
						fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
					}
					summary, err := test.ExecuteList(list, "--summary", "--iterations", "5")
					So(err, ShouldBeNil)

					So(summary.TotalRequests, ShouldEqual, float64(5))
				})
			})

			Convey("Workers", func() {

				Convey("Plan Usage", func() {

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

					So(err, ShouldBeNil)

					So(summary.TotalRequests, ShouldEqual, float64(5))
				})

				Convey("List Usage", func() {
					list := []string{
						fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
					}
					summary, err := test.ExecuteList(list, "--summary", "--workers", "5")
					So(err, ShouldBeNil)

					So(summary.TotalRequests, ShouldEqual, float64(5))
				})
			})

			Convey("Wait Time", func() {

				Convey("Plan Usage", func() {

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

					So(err, ShouldBeNil)

					So(math.Floor(summary.RunningTime.Seconds()), ShouldEqual, float64(5))
				})

				Convey("List Usage", func() {
					list := []string{
						fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
					}
					summary, err := test.ExecuteList(list, "--summary", "--wait-time", "1s", "--iterations", "5")
					So(err, ShouldBeNil)

					So(math.Floor(summary.RunningTime.Seconds()), ShouldEqual, float64(5))
				})
			})

			Convey("Duration", func() {

				Convey("Plan Usage", func() {

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

					So(err, ShouldBeNil)

					So(math.Floor(summary.RunningTime.Seconds()), ShouldEqual, float64(5))
				})

				Convey("List Usage", func() {
					list := []string{
						fmt.Sprintf(`%s -X GET'`, TestServer.CreateURL("/people")),
					}
					summary, err := test.ExecuteList(list, "--summary", "--duration", "5s")
					So(err, ShouldBeNil)

					So(math.Floor(summary.RunningTime.Seconds()), ShouldEqual, float64(5))
				})
			})
		})

		Convey("Halts execution if a payload input file is not found", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST -d '{"name":"talula"}'`, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST -d @missing-file.json`, URLForTestServer("/success")),
			}

			_, err := test.ExecuteList(list, "--summary")
			So(err, ShouldNotBeNil)

			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "Request body file not found: missing-file.json")
		})

		Convey("Error non-http url in the urls file causes a run time exception", func() {
			list := []string{
				fmt.Sprintf(`-Something`),
			}

			_, err := test.ExecuteList(list, "--summary")
			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, errormanager.LogMessageVaidURLs)
		})

		Convey("Issue - Should write out panics to a log file and not std out", func() {
			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				SetIterations(1).
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.IPanicAction().Build())

			output, err := test.ExecutePlanBuilder(planBuilder)
			So(err, ShouldNotBeNil)

			So(string(output), ShouldContainSubstring, "An unexpected error has occurred.  The error has been logged to /tmp/")

			//Ensure that the file which was generated contains the error which caused the panic
			r, _ := regexp.Compile(`/tmp/[\w\d-]+`)
			var location = r.FindString(string(output))
			So(location, ShouldNotEqual, "")
			data, err := ioutil.ReadFile(location)
			So(err, ShouldBeNil)
			So(string(data), ShouldContainSubstring, "IPanicAction has caused this panic")
		})
	})
}
