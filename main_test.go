package main

import (
	"io/ioutil"
	"math"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io"
	"net/http"

	"github.com/guzzlerio/rizo"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/global"
	"github.com/guzzlerio/corcel/infrastructure/inproc"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/test"
	. "github.com/guzzlerio/corcel/utils"
)

func URLForTestServer(path string) string {
	return TestServer.CreateURL(path)
}

var _ = Describe("Main", func() {

	AfterEach(func() {
		TestServer.Clear()
	})

	Describe("Support specified duraration for test", func() {
		It("Runs for 5 seconds", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			}

			output, err := SutExecuteApplication(list, config.Configuration{}.WithDuration("5s"))
			Expect(err).To(BeNil())

			var summary = statistics.CreateSummary(output)

			actual, _ := time.ParseDuration(summary.RunningTime)
			seconds := actual.Seconds()
			seconds = math.Floor(seconds)
			Expect(seconds).To(Equal(float64(5)))
		})
	})

	It("Support random selection of url from file", func() {
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
		Expect(err).To(BeNil())
		requestsSet1 := Requests(TestServer.Requests[:])
		TestServer.Clear()
		_, err = SutExecuteApplication(list, config.Configuration{
			Random: true,
		})
		Expect(err).To(BeNil())
		requestsSet2 := Requests(TestServer.Requests[:])

		Expect(ConcatRequestPaths(requestsSet1)).ToNot(Equal(ConcatRequestPaths(requestsSet2)))
	})

	for _, numberOfWorkers := range global.NumberOfWorkersToTest {
		func(workers int) {
			name := fmt.Sprintf("Support %v workers", workers)
			It(name, func() {
				list := []string{
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				}

				inproc.Throughput = 0
				inproc.ProcessEventsSubscribed = 0

				output, err := SutExecuteApplication(list, config.Configuration{
					Workers: workers,
				})
				Expect(err).To(BeNil())

				var summary = statistics.CreateSummary(output)
				//if summary.TotalRequests != float64(len(list)*workers) {
				fmt.Println(fmt.Sprintf(`
				Expected %v 
				Total Requests %v
				Process Events Subscribed %v
				`, float64(len(list)*workers),
					inproc.Throughput,
					inproc.ProcessEventsSubscribed))
				//}

				Expect(summary.TotalErrors).To(Equal(float64(0)))
				Expect(summary.TotalRequests).To(Equal(float64(len(list) * workers)))
			})
		}(numberOfWorkers)
	}

	for _, waitTime := range global.WaitTimeTests {
		It(fmt.Sprintf("Support wait time of %v between each execution in the list", waitTime), func() {
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
			Expect(err).To(BeNil())

			waitTimeValue, _ := time.ParseDuration(waitTime)
			expected := int64(len(list)) * int64(waitTimeValue)

			Expect(int64(duration)).To(BeNumerically(">=", int64(expected)))
		})
	}

	Describe("Generate statistics on throughput", func() {
		var list []string

		BeforeEach(func() {
			list = []string{
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			}
		})

		It("Records the availability", func() {
			count := 0
			TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
				count++
				if count%2 == 0 {
					w.WriteHeader(500)
				} else {
					w.WriteHeader(200)
				}
			}))

			output, err := SutExecuteApplication(list, config.Configuration{})
			Expect(err).To(BeNil())

			var summary = statistics.CreateSummary(output)

			Expect(summary.Availability).To(Equal(float64(60)))
		})

		for _, code := range global.ResponseCodes500 {
			It(fmt.Sprintf("Records error for HTTP %v response code range", code), func() {
				TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
					w.WriteHeader(code)
				}))

				output, err := SutExecuteApplication(list, config.Configuration{})
				Expect(err).To(BeNil())

				var summary = statistics.CreateSummary(output)

				Expect(summary.TotalErrors).To(Equal(float64(len(list))))
				Expect(summary.TotalRequests).To(Equal(float64(len(list))))
			})
		}

		It("Requests per second", func() {
			output, err := SutExecuteApplication(list, config.Configuration{})
			Expect(err).To(BeNil())

			var summary = statistics.CreateSummary(output)

			Expect(summary.Throughput).To(BeNumerically(">", 0))
			Expect(summary.TotalRequests).To(Equal(float64(len(list))))
		})
	})

	It("Generate statistics on timings", func() {
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

		output, err := SutExecuteApplication(list, config.Configuration{})
		Expect(err).To(BeNil())

		var summary = statistics.CreateSummary(output)

		Expect(summary.MaxResponseTime).To(BeNumerically(">", 0))
		Expect(summary.MeanResponseTime).To(BeNumerically(">", 0))
		Expect(summary.MinResponseTime).To(BeNumerically(">", 0))
	})

	It("Generate statistics of data from the execution", func() {
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

		output, err := SutExecuteApplication(list, config.Configuration{})
		Expect(err).To(BeNil())

		var summary = statistics.CreateSummary(output)

		Expect(summary.Bytes.TotalSent).To(BeNumerically(">", 0))

		Expect(summary.Bytes.TotalSent).To(BeNumerically(">", 0))
		Expect(summary.Bytes.TotalReceived).To(BeNumerically(">", 0))
	})

	Describe("Support sending data with http request", func() {
		for _, method := range global.HTTPMethodsWithRequestBody[:1] {

			It(fmt.Sprintf("in the body from a file for verb %s", method), func() {
				data := "@./list.txt"
				list := []string{fmt.Sprintf(`%s -X %s -d %s`, URLForTestServer("/A"), method, data)}
				_, err := SutExecuteApplication(list, config.Configuration{})
				Expect(err).To(BeNil())

				predicates := []rizo.HTTPRequestPredicate{}
				predicates = append(predicates, rizo.RequestWithPath("/A"))
				predicates = append(predicates, rizo.RequestWithMethod(method))

				bodyData, err := ioutil.ReadFile(data[1:])
				Expect(err).To(BeNil())
				predicates = append(predicates, rizo.RequestWithBody(string(bodyData)))
				Expect(TestServer.Find(predicates...)).To(Equal(true))
			})
		}

		It("in the querystring for verb GET", func() {
			method := "GET"
			data := "a=1&b=2&c=3"
			list := []string{fmt.Sprintf(`%s -X %s -d %s"`, URLForTestServer("/A"), method, data)}

			_, err := SutExecuteApplication(list, config.Configuration{})
			Expect(err).To(BeNil())

			predicates := []rizo.HTTPRequestPredicate{}
			predicates = append(predicates, rizo.RequestWithPath("/A"))
			predicates = append(predicates, rizo.RequestWithMethod(method))
			predicates = append(predicates, rizo.RequestWithQuerystring(data))
			Expect(TestServer.Find(predicates...)).To(Equal(true))
		})
	})

	for _, method := range global.SupportedHTTPMethods {
		It(fmt.Sprintf("Makes a http %s request", method), func() {
			list := []string{fmt.Sprintf(`%s -X %s`, URLForTestServer("/A"), method)}
			_, err := SutExecuteApplication(list, config.Configuration{})
			Expect(err).To(BeNil())
			Expect(TestServer.Find(rizo.RequestWithPath("/A"), rizo.RequestWithMethod(method))).To(Equal(true))
		})

		It(fmt.Sprintf("Makes a http %s request with http headers", method), func() {
			applicationJSON := "Content-Type:application/json"
			applicationSoapXML := "Accept:application/soap+xml"
			list := []string{fmt.Sprintf(`%s -X %s -H "%s" -H "%s"`, URLForTestServer("/A"), method, applicationJSON, applicationSoapXML)}

			_, err := SutExecuteApplication(list, config.Configuration{})
			Expect(err).To(BeNil())

			predicates := []rizo.HTTPRequestPredicate{}
			predicates = append(predicates, rizo.RequestWithPath("/A"))
			predicates = append(predicates, rizo.RequestWithMethod(method))
			predicates = append(predicates, rizo.RequestWithHeader("Content-Type", "application/json"))
			predicates = append(predicates, rizo.RequestWithHeader("Accept", "application/soap+xml"))
			Expect(TestServer.Find(predicates...)).To(Equal(true))
		})
	}

	It("Makes a http request with a custom user agent", func() {
		userAgent := "Mozilla/5.0 (iPhone; U; CPU iPhone OS 5_1_1 like Mac OS X; en) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3"

		method := "POST"
		list := []string{fmt.Sprintf(`%s -X %s -A "%s"`, URLForTestServer("/A"), method, userAgent)}
		_, err := SutExecuteApplication(list, config.Configuration{})
		Expect(err).To(BeNil())

		predicates := []rizo.HTTPRequestPredicate{}
		predicates = append(predicates, rizo.RequestWithPath("/A"))
		predicates = append(predicates, rizo.RequestWithMethod(method))
		predicates = append(predicates, rizo.RequestWithHeader("User-Agent", userAgent))
		Expect(TestServer.Find(predicates...)).To(Equal(true))
	})

	It("Makes a http get request to each url in a file", func() {
		list := []string{
			URLForTestServer("/A"),
			URLForTestServer("/B"),
			URLForTestServer("/C"),
		}

		_, err := SutExecuteApplication(list, config.Configuration{})
		Expect(err).To(BeNil())

		Expect(TestServer.Find(rizo.RequestWithPath("/A"))).To(Equal(true))
		Expect(TestServer.Find(rizo.RequestWithPath("/B"))).To(Equal(true))
		Expect(TestServer.Find(rizo.RequestWithPath("/C"))).To(Equal(true))
	})
})

func SutExecuteApplication(list []string, configuration config.Configuration) (statistics.AggregatorSnapShot, error) {
	output, err := test.ExecuteListForApplication(list, configuration)
	return output, err
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
