package http_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	"github.com/guzzlerio/rizo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExecutionPlanHttpRequest(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("ExecutionPlanHttpRequest", t, func() {
		func() {
			TestServer.Clear()
			factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				w.Header().Add("X-BOOM", "1")
				w.WriteHeader(http.StatusOK)
			})

			TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
		}()

		Convey("Adds the response headers to the context", func() {

			plan := fmt.Sprintf(`---
workers: 1
jobs:
- name: "Job 1"
  steps:
  - name: "Step 1"
    action:
      type: HttpRequest
      headers:
        key: 1
      method: GET
      url: %s
    assertions:
    - type: ExactAssertion
      key: urn:http:response:headers:x-boom
      expected: "1"`, TestServer.CreateURL("/people"))

			summary, err := test.ExecutePlanFromDataForApplication(plan)
			So(err, ShouldBeNil)
			So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
		})

		Convey("Supplies a payload to the HTTP Request", func() {
			planBuilder := yaml.NewPlanBuilder()

			path := "/people"
			body := "Zee Body"

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

			_, err := test.ExecutePlanBuilder(planBuilder)
			So(err, ShouldBeNil)
			So(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithBody(body)), ShouldEqual, true)
		})

		Convey("Supplies a header which is an int", func() {
			plan := fmt.Sprintf(`---
iterations: 0
random: false
workers: 1
waitTime: 0s
duration: 0s
jobs:
- name: ""
  steps:
  - name: ""
    action:
      body: Zee Body
      headers:
        key: 1
      method: GET
      type: HttpRequest
      url: %s
`, TestServer.CreateURL("/people"))

			_, err := test.ExecutePlanFromData(plan)
			So(err, ShouldBeNil)
			So(TestServer.Find(rizo.RequestWithHeader("key", "1")), ShouldEqual, true)
		})

		Convey("Supplies a payload as a file reference to the HTTP Request", func() {
			content := []byte("temporary file's content")
			dir, err := ioutil.TempDir("", "ExecutionPlanHttpRequest")
			if err != nil {
				panic(err)
			}

			defer func() {
				err := os.RemoveAll(dir) // clean up
				if err != nil {
					panic(err)
				}
			}()

			tmpfn := filepath.Join(dir, "tmpfile")
			if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
				panic(err)
			}

			path := "/people"
			body := fmt.Sprintf("@%s", tmpfn)

			planBuilder := yaml.NewPlanBuilder()

			planBuilder.
				CreateJob().
				CreateStep().
				ToExecuteAction(planBuilder.HTTPAction().URL(TestServer.CreateURL(path)).Body(body).Build())

			_, err = test.ExecutePlanBuilder(planBuilder)
			So(err, ShouldBeNil)
			So(TestServer.Find(rizo.RequestWithPath(path), rizo.RequestWithBody(string(content))), ShouldEqual, true)
		})

		SkipConvey("Returns an error when the file path does not exist", func() {})
	})
}
