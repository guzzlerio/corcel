package main

import (
	"io/ioutil"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/guzzlerio/rizo"

	"ci.guzzler.io/guzzler/corcel/errormanager"
	"ci.guzzler.io/guzzler/corcel/global"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/statistics"
	. "ci.guzzler.io/guzzler/corcel/utils"
)

func URLForTestServer(path string) string {
	return TestServer.CreateURL(path)
}

var _ = Describe("Main", func() {
	BeforeEach(func() {
		err := os.Remove("./output.yml")
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	})

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

			SutExecute(list, "--duration", "5s")

			var executionOutput statistics.AggregatorSnapShot
			UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

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

		SutExecute(list, "--random")
		requestsSet1 := Requests(TestServer.Requests[:])
		TestServer.Clear()
		SutExecute(list, "--random")
		requestsSet2 := Requests(TestServer.Requests[:])

		Expect(ConcatRequestPaths(requestsSet1)).ToNot(Equal(ConcatRequestPaths(requestsSet2)))
	})

	for _, numberOfWorkers := range global.NumberOfWorkersToTest {
		name := fmt.Sprintf("Support %v workers", numberOfWorkers)
		It(name, func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			}

			SutExecute(list, "--workers", strconv.Itoa(numberOfWorkers))

			var executionOutput statistics.AggregatorSnapShot
			UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.TotalRequests).To(Equal(float64(len(list) * numberOfWorkers)))
			Expect(summary.TotalErrors).To(Equal(float64(0)))

		})
	}

	for _, waitTime := range global.WaitTimeTests {
		It(fmt.Sprintf("Support wait time of %v between each execution in the list", waitTime), func() {
			waitTimeTolerance := 0.75

			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			}
			start := time.Now()
			SutExecute(list, "--wait-time", waitTime)
			duration := time.Since(start)

			waitTimeValue, _ := time.ParseDuration(waitTime)
			expected := int64(len(list)) * int64(waitTimeValue)
			maximum := float64(expected) * (1 + waitTimeTolerance)

			Expect(int64(duration)).To(BeNumerically(">=", int64(expected)))
			Expect(int64(duration)).To(BeNumerically("<", int64(maximum)))
		})
	}

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

		output := SutExecute(list, "--summary")

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Running Time: %v", summary.RunningTime)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Throughput: %.0f req/s", summary.Throughput)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Total Requests: %v", summary.TotalRequests)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Number of Errors: %v", summary.TotalErrors)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Availability: %v.0000%%", summary.Availability)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Bytes Sent: %v", summary.TotalBytesSent)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Bytes Received: %v", summary.TotalBytesReceived)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Mean Response Time: %.4f", summary.MeanResponseTime)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Min Response Time: %.4f ms", summary.MinResponseTime)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Max Response Time: %.4f ms", summary.MaxResponseTime)))
	})

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

			SutExecute(list)

			var executionOutput statistics.AggregatorSnapShot
			UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

			Expect(summary.Availability).To(Equal(float64(60)))
		})

		for _, code := range global.ResponseCodes500 {
			It(fmt.Sprintf("Records error for HTTP %v response code range", code), func() {
				TestServer.Use(rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
					w.WriteHeader(code)
				}))

				SutExecute(list)

				var executionOutput statistics.AggregatorSnapShot
				UnmarshalYamlFromFile("./output.yml", &executionOutput)
				var summary = statistics.CreateSummary(executionOutput)

				Expect(summary.TotalErrors).To(Equal(float64(len(list))))
				Expect(summary.TotalRequests).To(Equal(float64(len(list))))
			})
		}

		It("Requests per second", func() {
			SutExecute(list)

			var executionOutput statistics.AggregatorSnapShot
			UnmarshalYamlFromFile("./output.yml", &executionOutput)
			var summary = statistics.CreateSummary(executionOutput)

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

		SutExecute(list)

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(summary.MaxResponseTime).To(BeNumerically(">", 0))
		Expect(summary.MeanResponseTime).To(BeNumerically(">", 0))
		Expect(summary.MinResponseTime).To(BeNumerically(">", 0))
	})

	It("Halts execution if a payload input file is not found", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST -d '{"name":"talula"}'`, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X POST -d @missing-file.json`, URLForTestServer("/success")),
		}

		output, _ := InvokeCorcel(list)

		Expect(string(output)).To(ContainSubstring("Request body file not found: missing-file.json"))
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

		SutExecute(list)

		Expect(PathExists("./output.yml")).To(Equal(true))

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(summary.TotalBytesSent).To(BeNumerically(">", 0))

		Expect(summary.TotalBytesSent).To(BeNumerically(">", 0))
		Expect(summary.TotalBytesReceived).To(BeNumerically(">", 0))
	})

	Describe("Support sending data with http request", func() {
		for _, method := range global.HTTPMethodsWithRequestBody[:1] {
			PIt(fmt.Sprintf("in the body for verb %s", method), func() {
				data := "a=1&b=2&c=3"
				list := []string{fmt.Sprintf(`%s -X %s -d %s`, URLForTestServer("/A"), method, data)}
				SutExecute(list)

				predicates := []rizo.HTTPRequestPredicate{}
				predicates = append(predicates, rizo.RequestWithPath("/A"))
				predicates = append(predicates, rizo.RequestWithMethod(method))
				predicates = append(predicates, rizo.RequestWithBody(data))
				Expect(TestServer.Find(predicates...)).To(Equal(true))
			})

			It(fmt.Sprintf("in the body from a file for verb %s", method), func() {
				data := "@./list.txt"
				list := []string{fmt.Sprintf(`%s -X %s -d %s`, URLForTestServer("/A"), method, data)}
				SutExecute(list)

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
			SutExecute(list)

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
			SutExecute(list)
			Expect(TestServer.Find(rizo.RequestWithPath("/A"), rizo.RequestWithMethod(method))).To(Equal(true))
		})

		It(fmt.Sprintf("Makes a http %s request with http headers", method), func() {
			applicationJSON := "Content-Type:application/json"
			applicationSoapXML := "Accept:application/soap+xml"
			list := []string{fmt.Sprintf(`%s -X %s -H "%s" -H "%s"`, URLForTestServer("/A"), method, applicationJSON, applicationSoapXML)}
			SutExecute(list)

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
		SutExecute(list)

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

		SutExecute(list)

		Expect(TestServer.Find(rizo.RequestWithPath("/A"))).To(Equal(true))
		Expect(TestServer.Find(rizo.RequestWithPath("/B"))).To(Equal(true))
		Expect(TestServer.Find(rizo.RequestWithPath("/C"))).To(Equal(true))
	})
})

func InvokeCorcel(list []string, args ...string) ([]byte, error) {
	exePath, exeErr := filepath.Abs("./corcel")
	check(exeErr)
	file := CreateFileFromLines(list)
	defer func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	}()
	cmd := exec.Command(exePath, append(append([]string{"run", "--progress", "none"}, args...), file.Name())...)
	output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	if len(output) > 0 {
		logger.Log.Println(fmt.Sprintf("%s", output))
	}
	return output, err
}

func SutExecute(list []string, args ...string) []byte {
	output, err := InvokeCorcel(list, args...)
	if err != nil {
		Fail(string(output))
	}
	return output
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
