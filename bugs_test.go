package main

import (
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/guzzlerio/rizo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/guzzlerio/corcel/config"
)

func TestBugs_replication(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("Bugs replication", t, func() {

		defer func() {
			TestServer.Clear()
		}()

		Convey("Error when running a simple run with POST and PUT #18", func() {
			numberOfWorkers := 2
			list := []string{
				fmt.Sprintf(`%s -X POST -d '{"name": "bob"}' -H "Content-type: application/json"`, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X PUT -d '{"id": 1,"name": "bob junior"}' -H "Content-type: application/json"`, URLForTestServer("/success")),
				fmt.Sprintf(`%s?id=1 -X GET -H "Content-type: application/json"`, URLForTestServer("/success")),
				fmt.Sprintf(`%s?id=1 -X DELETE -H "Content-type: application/json"`, URLForTestServer("/success")),
			}

			summary, err := SutExecuteApplication(list[:1], config.Configuration{
				Random:  true,
				Summary: true,
				Workers: numberOfWorkers,
			})
			So(err, ShouldBeNil)
			So(summary.TotalRequests, ShouldEqual, float64(2))
		})

		Convey("Issue #49 - Corcel not cancelling on-going requests once the test is due to finish", func() {
			TestServer.Clear()
			factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				time.Sleep(2 * time.Second)
				w.WriteHeader(http.StatusOK)
			})
			TestServer.Use(factory)
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/something")),
			}

			summary, err := SutExecuteApplication(list, config.Configuration{}.WithDuration("1s"))
			So(err, ShouldBeNil)
			runningTime := summary.RunningTime
			So(math.Floor(runningTime.Seconds()), ShouldEqual, float64(1))
		})
	})
}
