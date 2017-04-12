package main

import (
	"io/ioutil"
	"math"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	"io"
	"net/http"

	"github.com/guzzlerio/rizo"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/test"
	. "github.com/guzzlerio/corcel/utils"
)

func URLForTestServer(path string) string {
	return TestServer.CreateURL(path)
}

func TestMain(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("Main", t, func() {

		defer func() {
			TestServer.Clear()
		}()

		Convey("Support specified duraration for test", func() {
			Convey("Runs for 5 seconds", func() {
				list := []string{
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				}

				summary, err := SutExecuteApplication(list, config.Configuration{}.WithDuration("5s"))
				So(err, ShouldBeNil)
				actual := summary.RunningTime
				seconds := actual.Seconds()
				seconds = math.Floor(seconds)
				So(seconds, ShouldEqual, float64(5))
			})
		})

		Convey("Support random selection of url from file", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/1")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/2")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/3")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/4")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/5")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/6")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/7")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/8")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/9")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/10")),
			}

			_, err := SutExecuteApplication(list, config.Configuration{
				Random: true,
			})
			So(err, ShouldBeNil)
			requestsSet1 := Requests(TestServer.Requests[:])
			TestServer.Clear()
			_, err = SutExecuteApplication(list, config.Configuration{
				Random: true,
			})
			So(err, ShouldBeNil)
			requestsSet2 := Requests(TestServer.Requests[:])

			So(ConcatRequestPaths(requestsSet1), ShouldNotResemble, ConcatRequestPaths(requestsSet2))
		})

		for _, numberOfWorkers := range global.NumberOfWorkersToTest {
			func(workers int) {
				name := fmt.Sprintf("Support %v workers", workers)
				Convey(name, func() {
					list := []string{
						fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
						fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
						fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
						fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
						fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
						fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					}

					summary, err := SutExecuteApplication(list, config.Configuration{
						Workers: workers,
					})
					So(err, ShouldBeNil)
					So(summary.TotalErrors, ShouldEqual, float64(0))
					So(summary.TotalRequests, ShouldEqual, float64(len(list)*workers))
				})
			}(numberOfWorkers)
		}

		for _, waitTime := range global.WaitTimeTests {
			Convey(fmt.Sprintf("Support wait time of %v between each execution in the list", waitTime), func() {
				list := []string{
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				}
				start := time.Now()
				_, err := SutExecuteApplication(list, config.Configuration{}.WithWaitTime(waitTime))
				duration := time.Since(start)
				So(err, ShouldBeNil)

				waitTimeValue, _ := time.ParseDuration(waitTime)
				expected := int64(len(list)) * int64(waitTimeValue)

				So(int64(duration), ShouldBeGreaterThanOrEqualTo, int64(expected))
			})
		}

		Convey("Generate statistics on throughput", func() {
			var list []string

			func() {
				list = []string{
					fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
					fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
					fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
					fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
					fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				}
			}()

			Convey("Records the availability", func() {
				count := 0
				TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
					count++
					if count%2 == 0 {
						w.WriteHeader(500)
					} else {
						w.WriteHeader(200)
					}
				}))

				summary, err := SutExecuteApplication(list, config.Configuration{})
				So(err, ShouldBeNil)
				So(summary.Availability, ShouldEqual, float64(60))
			})

			for _, code := range global.ResponseCodes500 {
				Convey(fmt.Sprintf("Records error for HTTP %v response code range", code), func() {
					TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
						w.WriteHeader(code)
					}))

					summary, err := SutExecuteApplication(list, config.Configuration{})
					So(err, ShouldBeNil)
					So(summary.TotalErrors, ShouldEqual, float64(len(list)))
					So(summary.TotalRequests, ShouldEqual, float64(len(list)))
				})
			}

			Convey("Requests per second", func() {
				summary, err := SutExecuteApplication(list, config.Configuration{})
				So(err, ShouldBeNil)
				So(summary.Throughput, ShouldBeGreaterThan, 0)
				So(summary.TotalRequests, ShouldEqual, float64(len(list)))
			})
		})

		Convey("Generate statistics on timings", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			}

			count := 1
			TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				time.Sleep(time.Duration(count) * time.Millisecond)
				count++
			}))

			summary, err := SutExecuteApplication(list, config.Configuration{})
			So(err, ShouldBeNil)
			So(summary.ResponseTime.Max, ShouldBeGreaterThan, 0)
			So(summary.ResponseTime.Mean, ShouldBeGreaterThan, 0)
			So(summary.ResponseTime.Min, ShouldBeGreaterThan, 0)
		})

		Convey("Generate statistics of data from the execution", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X PUT -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X DELETE -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X GET`, URLForTestServer("/A")),
			}

			responseBody := "-"
			TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				_, err := io.WriteString(w, fmt.Sprintf("%s", responseBody))
				check(err)
				responseBody = responseBody + "-"
			}))

			summary, err := SutExecuteApplication(list, config.Configuration{})
			So(err, ShouldBeNil)
			So(summary.Bytes.Sent.Total, ShouldBeGreaterThan, 0)
			So(summary.Bytes.Received.Total, ShouldBeGreaterThan, 0)
		})

		Convey("Support sending data with http request", func() {
			for _, method := range global.HTTPMethodsWithRequestBody[:1] {

				Convey(fmt.Sprintf("in the body from a file for verb %s", method), func() {
					data := "@./list.txt"
					list := []string{fmt.Sprintf(`%s -X %s -d %s`, URLForTestServer("/A"), method, data)}
					_, err := SutExecuteApplication(list, config.Configuration{})
					So(err, ShouldBeNil)

					predicates := []rizo.HTTPRequestPredicate{}
					predicates = append(predicates, rizo.RequestWithPath("/A"))
					predicates = append(predicates, rizo.RequestWithMethod(method))

					bodyData, err := ioutil.ReadFile(data[1:])
					So(err, ShouldBeNil)
					predicates = append(predicates, rizo.RequestWithBody(string(bodyData)))
					So(TestServer.Find(predicates...), ShouldEqual, true)
				})
			}

			Convey("in the querystring for verb GET", func() {
				method := "GET"
				data := "a=1&b=2&c=3"
				list := []string{fmt.Sprintf(`%s -X %s -d %s"`, URLForTestServer("/A"), method, data)}

				_, err := SutExecuteApplication(list, config.Configuration{})
				So(err, ShouldBeNil)

				predicates := []rizo.HTTPRequestPredicate{}
				predicates = append(predicates, rizo.RequestWithPath("/A"))
				predicates = append(predicates, rizo.RequestWithMethod(method))
				predicates = append(predicates, rizo.RequestWithQuerystring(data))
				So(TestServer.Find(predicates...), ShouldEqual, true)
			})
		})

		for _, method := range global.SupportedHTTPMethods {
			Convey(fmt.Sprintf("Makes a http %s request", method), func() {
				list := []string{fmt.Sprintf(`%s -X %s`, URLForTestServer("/A"), method)}
				_, err := SutExecuteApplication(list, config.Configuration{})
				So(err, ShouldBeNil)
				So(TestServer.Find(rizo.RequestWithPath("/A"), rizo.RequestWithMethod(method)), ShouldEqual, true)
			})

			Convey(fmt.Sprintf("Makes a http %s request with http headers", method), func() {
				applicationJSON := "Content-Type:application/json"
				applicationSoapXML := "Accept:application/soap+xml"
				list := []string{fmt.Sprintf(`%s -X %s -H "%s" -H "%s"`, URLForTestServer("/A"), method, applicationJSON, applicationSoapXML)}

				_, err := SutExecuteApplication(list, config.Configuration{})
				So(err, ShouldBeNil)

				predicates := []rizo.HTTPRequestPredicate{}
				predicates = append(predicates, rizo.RequestWithPath("/A"))
				predicates = append(predicates, rizo.RequestWithMethod(method))
				predicates = append(predicates, rizo.RequestWithHeader("Content-Type", "application/json"))
				predicates = append(predicates, rizo.RequestWithHeader("Accept", "application/soap+xml"))
				So(TestServer.Find(predicates...), ShouldEqual, true)
			})
		}

		Convey("Makes a http request with a custom user agent", func() {
			userAgent := "Mozilla/5.0 (iPhone; U; CPU iPhone OS 5_1_1 like Mac OS X; en) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3"

			method := "POST"
			list := []string{fmt.Sprintf(`%s -X %s -A "%s"`, URLForTestServer("/A"), method, userAgent)}
			_, err := SutExecuteApplication(list, config.Configuration{})
			So(err, ShouldBeNil)

			predicates := []rizo.HTTPRequestPredicate{}
			predicates = append(predicates, rizo.RequestWithPath("/A"))
			predicates = append(predicates, rizo.RequestWithMethod(method))
			predicates = append(predicates, rizo.RequestWithHeader("User-Agent", userAgent))
			So(TestServer.Find(predicates...), ShouldEqual, true)
		})

		Convey("Makes a http get request to each url in a file", func() {
			list := []string{
				URLForTestServer("/A"),
				URLForTestServer("/B"),
				URLForTestServer("/C"),
			}

			_, err := SutExecuteApplication(list, config.Configuration{})
			So(err, ShouldBeNil)

			So(TestServer.Find(rizo.RequestWithPath("/A")), ShouldEqual, true)
			So(TestServer.Find(rizo.RequestWithPath("/B")), ShouldEqual, true)
			So(TestServer.Find(rizo.RequestWithPath("/C")), ShouldEqual, true)
		})
	})
}

func SutExecuteApplication(list []string, configuration config.Configuration) (core.ExecutionSummary, error) {
	return test.ExecuteListForApplication(list, configuration)
}

func Requests(recordedRequests []rizo.RecordedRequest) (result []*http.Request) {
	for _, recordedRequest := range recordedRequests {
		result = append(result, recordedRequest.Request)
	}
	return
}

func check(err error) {
	if err != nil {
		errormanager.Log(err)
	}
}
