package main

import (
	"io"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var (
	SUPPORTED_HTTP_METHODS         = []string{"GET", "POST", "PUT", "DELETE"}
	HTTP_METHODS_WITH_REQUEST_BODY = []string{"POST", "PUT", "DELETE"}
	TestServer                     *RequestRecordingServer
	TEST_PORT                      = 8000
	RESPONSE_CODES_400             = []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418}
	RESPONSE_CODES_500             = []int{500, 501, 502, 503, 504, 505}
	WAIT_TIME_TESTS                = []string{"1ms", "2ms", "4ms", "8ms", "16ms", "32ms", "64ms", "128ms"}
	NUMBER_OF_WORKERS_TO_TEST      = []int{1, 2, 4, 8, 16, 32, 64, 128, 256}
)

func UrlForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TEST_PORT, path)
}

var _ = BeforeSuite(func() {
	ConfigureLogging()
	TestServer = CreateRequestRecordingServer(TEST_PORT)
	TestServer.Start()
})

var _ = AfterSuite(func() {
	TestServer.Stop()
})

func SutExecute(list []string, args ...string) []byte {
	exePath, err := filepath.Abs("./corcel")
	check(err)
	file := CreateFileFromLines(list)
	defer os.Remove(file.Name())
	cmd := exec.Command(exePath, append([]string{"-f", file.Name()}, args...)...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		Log.Println(string(output))
	}
	Expect(err).To(BeNil())
	return output
}

func Requests(recordedRequests []RecordedRequest) (result []*http.Request) {
	for _, recordedRequest := range recordedRequests {
		result = append(result, recordedRequest.request)
	}
	return
}

func ConcatRequestPaths(requests []*http.Request) string {
	result := ""
	for _, request := range requests {
		result = result + request.URL.Path
	}
	return result
}

var _ = Describe("Main", func() {

	var (
		exePath string
		err     error
	)

	BeforeEach(func() {
		os.Remove("./output.yml")
		exePath, err = filepath.Abs("./corcel")
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	Describe("Support specified duraration for test", func() {
		It("Runs for 10 seconds", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
			}

			SutExecute(list, "--duration", "5s")

			var executionOutput ExecutionOutput
			UnmarshalYamlFromFile("./output.yml", &executionOutput)

			Expect(int64(executionOutput.Summary.RunningTime)).To(BeNumerically(">=", int64(5000)))
			Expect(int64(executionOutput.Summary.RunningTime)).To(BeNumerically("<", int64(6000)))
		})
	})

	It("Support random selection of url from file", func() {

		list := []string{
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/1")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/2")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/3")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/4")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/5")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/6")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/7")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/8")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/9")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/10")),
		}

		SutExecute(list, "--random")
		requestsSet1 := Requests(TestServer.requests[:])
		TestServer.Clear()
		SutExecute(list, "--random")
		requestsSet2 := Requests(TestServer.requests[:])

		Expect(ConcatRequestPaths(requestsSet1)).ToNot(Equal(ConcatRequestPaths(requestsSet2)))

	})

	for _, numberOfWorkers := range NUMBER_OF_WORKERS_TO_TEST {
		It(fmt.Sprintf("Support %v workers", numberOfWorkers), func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
			}

			SutExecute(list, "--workers", strconv.Itoa(numberOfWorkers))

			var executionOutput ExecutionOutput
			UnmarshalYamlFromFile("./output.yml", &executionOutput)

			Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(len(list) * numberOfWorkers)))
			Expect(executionOutput.Summary.Requests.Errors).To(Equal(int64(0)))

		})
	}

	for _, waitTime := range WAIT_TIME_TESTS {
		It(fmt.Sprintf("Support wait time of %v between each execution in the list", waitTime), func() {
			waitTimeTolerance := 0.25

			list := []string{
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
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
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/error")),
			fmt.Sprintf(`%s -X POST `, UrlForTestServer("/success")),
		}

		TestServer.Use(HttpResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(500)
		})).For(RequestWithPath("/error"))

		output := SutExecute(list, "--summary")

		var executionOutput ExecutionOutput
		UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Running Time: %v s", executionOutput.Summary.RunningTime/1000)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Throughput: %v req/s", int(executionOutput.Summary.Requests.Rate))))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Total Requests: %v", executionOutput.Summary.Requests.Total)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Number of Errors: %v", executionOutput.Summary.Requests.Errors)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Availability: %v%%", executionOutput.Summary.Requests.Availability*100)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Bytes Sent: %v", executionOutput.Summary.Bytes.Sent.Sum)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Bytes Received: %v", executionOutput.Summary.Bytes.Received.Sum)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Mean Response Time: %.4v", executionOutput.Summary.ResponseTime.Mean)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Min Response Time: %v ms", executionOutput.Summary.ResponseTime.Min)))
		Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Max Response Time: %v ms", executionOutput.Summary.ResponseTime.Max)))
	})

	Describe("Generate statistics on throughput", func() {
		var list []string

		BeforeEach(func() {
			list = []string{
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			}
		})

		It("Records the availability", func() {
			count := 0
			TestServer.Use(HttpResponseFactory(func(w http.ResponseWriter) {
				count++
				if count%2 == 0 {
					w.WriteHeader(500)
				} else {
					w.WriteHeader(200)
				}
			}))

			SutExecute(list)

			var executionOutput ExecutionOutput
			UnmarshalYamlFromFile("./output.yml", &executionOutput)

			expectedAvailability := 1 - (float64(executionOutput.Summary.Requests.Errors) / float64(executionOutput.Summary.Requests.Total))
			Expect(executionOutput.Summary.Requests.Availability).To(Equal(expectedAvailability))
		})

		for _, code := range append(RESPONSE_CODES_500, RESPONSE_CODES_400...) {
			It(fmt.Sprintf("Records error for HTTP %v response code range", code), func() {
				TestServer.Use(HttpResponseFactory(func(w http.ResponseWriter) {
					w.WriteHeader(code)
				}))

				SutExecute(list)

				var executionOutput ExecutionOutput
				UnmarshalYamlFromFile("./output.yml", &executionOutput)

				Expect(executionOutput.Summary.Requests.Errors).To(Equal(int64(len(list))))
				Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(len(list))))
			})
		}

		It("Requests per second", func() {
			SutExecute(list)

			var executionOutput ExecutionOutput
			UnmarshalYamlFromFile("./output.yml", &executionOutput)
			Expect(executionOutput.Summary.Requests.Rate).To(BeNumerically(">", 0))
			Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(len(list))))
		})
	})

	It("Generate statistics on timings", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
		}

		count := 1
		TestServer.Use(HttpResponseFactory(func(w http.ResponseWriter) {
			time.Sleep(time.Duration(count) * time.Millisecond)
			count++
		}))

		SutExecute(list)

		var executionOutput ExecutionOutput

		UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(executionOutput.Summary.ResponseTime.Sum).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.Max).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.Mean).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.Min).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.P50).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.P75).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.P95).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.P99).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.StdDev).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.ResponseTime.Var).To(BeNumerically(">", 0))

		Expect(executionOutput.Summary.RunningTime).To(BeNumerically(">", 0))
	})

	It("Generate statistics of data from the execution", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X PUT -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X DELETE -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X GET`, UrlForTestServer("/A")),
		}

		responseBody := "-"
		TestServer.Use(HttpResponseFactory(func(w http.ResponseWriter) {
			io.WriteString(w, fmt.Sprintf("%s", responseBody))
			responseBody = responseBody + "-"
		}))

		SutExecute(list)

		Expect(PathExists("./output.yml")).To(Equal(true))

		var executionOutput ExecutionOutput

		UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(executionOutput.Summary.Bytes.Sent.Sum).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.Max).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.Mean).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.Min).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.P50).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.P75).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.P95).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.P99).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.StdDev).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.Var).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Sent.Rate).To(BeNumerically(">", 0))

		Expect(executionOutput.Summary.Bytes.Received.Sum).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.Max).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.Mean).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.Min).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.P50).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.P75).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.P95).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.P99).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.StdDev).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.Var).To(BeNumerically(">", 0))
		Expect(executionOutput.Summary.Bytes.Received.Rate).To(BeNumerically(">", 0))
	})

	Describe("Support sending data with http request", func() {
		for _, method := range HTTP_METHODS_WITH_REQUEST_BODY {
			It(fmt.Sprintf("in the body for verb %s", method), func() {
				data := "a=1&b=2&c=3"
				list := []string{fmt.Sprintf(`%s -X %s -d %s`, UrlForTestServer("/A"), method, data)}
				SutExecute(list)

				predicates := []HttpRequestPredicate{}
				predicates = append(predicates, RequestWithPath("/A"))
				predicates = append(predicates, RequestWithMethod(method))
				predicates = append(predicates, RequestWithBody(data))
				Expect(TestServer.Find(predicates...)).To(Equal(true))
			})
		}

		It("in the querystring for verb GET", func() {
			method := "GET"
			data := "a=1&b=2&c=3"
			list := []string{fmt.Sprintf(`%s -X %s -d %s"`, UrlForTestServer("/A"), method, data)}
			SutExecute(list)

			predicates := []HttpRequestPredicate{}
			predicates = append(predicates, RequestWithPath("/A"))
			predicates = append(predicates, RequestWithMethod(method))
			predicates = append(predicates, RequestWithQuerystring(data))
			Expect(TestServer.Find(predicates...)).To(Equal(true))
		})
	})

	for _, method := range SUPPORTED_HTTP_METHODS {
		It(fmt.Sprintf("Makes a http %s request with http headers", method), func() {
			applicationJson := "Content-Type:application/json"
			applicationSoapXml := "Accept:application/soap+xml"
			list := []string{fmt.Sprintf(`%s -X %s -H "%s" -H "%s"`, UrlForTestServer("/A"), method, applicationJson, applicationSoapXml)}
			SutExecute(list)

			predicates := []HttpRequestPredicate{}
			predicates = append(predicates, RequestWithPath("/A"))
			predicates = append(predicates, RequestWithMethod(method))
			predicates = append(predicates, RequestWithHeader("Content-Type", "application/json"))
			predicates = append(predicates, RequestWithHeader("Accept", "application/soap+xml"))
			Expect(TestServer.Find(predicates...)).To(Equal(true))
		})
	}

	for _, method := range SUPPORTED_HTTP_METHODS {
		It(fmt.Sprintf("Makes a http %s request", method), func() {
			list := []string{fmt.Sprintf(`%s -X %s`, UrlForTestServer("/A"), method)}
			SutExecute(list)
			Expect(TestServer.Find(RequestWithPath("/A"), RequestWithMethod(method))).To(Equal(true))
		})
	}

	It("Makes a http get request to each url in a file", func() {
		list := []string{
			UrlForTestServer("/A"),
			UrlForTestServer("/B"),
			UrlForTestServer("/C"),
		}

		SutExecute(list)

		Expect(TestServer.Find(RequestWithPath("/A"))).To(Equal(true))
		Expect(TestServer.Find(RequestWithPath("/B"))).To(Equal(true))
		Expect(TestServer.Find(RequestWithPath("/C"))).To(Equal(true))
	})
})
