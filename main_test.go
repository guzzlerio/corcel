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
	SupportedHTTPMethods       = []string{"GET", "POST", "PUT", "DELETE"}
	HTTPMethodsWithRequestBody = []string{"POST", "PUT", "DELETE"}
	TestServer                 *RequestRecordingServer
	TestPort                   = 8000
	ResponseCodes400           = []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418}
	ResponseCodes500           = []int{500, 501, 502, 503, 504, 505}
	WaitTimeTests              = []string{"1ms", "2ms", "4ms", "8ms", "16ms", "32ms", "64ms", "128ms"}
	NumberOfWorkersToTest      = []int{1, 2, 4, 8, 16, 32, 64, 128, 256}
)

func URLForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TestPort, path)
}

var _ = BeforeSuite(func() {
	ConfigureLogging()
	TestServer = CreateRequestRecordingServer(TestPort)
	TestServer.Start()
})

var _ = AfterSuite(func() {
	TestServer.Stop()
})

func InvokeCorcel(list []string, args ...string) ([]byte, error) {
	exePath, err := filepath.Abs("./corcel")
	check(err)
	configFileReader = func(path string) ([]byte, error) {
		return []byte(""), nil
	}
	file := CreateFileFromLines(list)
	defer os.Remove(file.Name())
	cmd := exec.Command(exePath, append(args, file.Name())...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		Log.Println(fmt.Sprintf("%s", output))
	}
	return output, err
}

func SutExecute(list []string, args ...string) []byte {
	output, err := InvokeCorcel(list, args...)
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

	BeforeEach(func() {
		os.Remove("./output.yml")
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	Describe("Support specified duraration for test", func() {
		It("Runs for 10 seconds", func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
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
		requestsSet1 := Requests(TestServer.requests[:])
		TestServer.Clear()
		SutExecute(list, "--random")
		requestsSet2 := Requests(TestServer.requests[:])

		Expect(ConcatRequestPaths(requestsSet1)).ToNot(Equal(ConcatRequestPaths(requestsSet2)))
	})

	for _, numberOfWorkers := range NumberOfWorkersToTest {
		It(fmt.Sprintf("Support %v workers", numberOfWorkers), func() {
			list := []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/success")),
			}

			SutExecute(list, "--workers", strconv.Itoa(numberOfWorkers))

			var executionOutput ExecutionOutput
			UnmarshalYamlFromFile("./output.yml", &executionOutput)

			Expect(executionOutput.Summary.Requests.Total).To(Equal(int64(len(list) * numberOfWorkers)))
			Expect(executionOutput.Summary.Requests.Errors).To(Equal(int64(0)))

		})
	}

	for _, waitTime := range WaitTimeTests {
		It(fmt.Sprintf("Support wait time of %v between each execution in the list", waitTime), func() {
			waitTimeTolerance := 0.25

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

		TestServer.Use(HTTPResponseFactory(func(w http.ResponseWriter) {
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
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
				fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			}
		})

		It("Records the availability", func() {
			count := 0
			TestServer.Use(HTTPResponseFactory(func(w http.ResponseWriter) {
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

		for _, code := range append(ResponseCodes500, ResponseCodes400...) {
			It(fmt.Sprintf("Records error for HTTP %v response code range", code), func() {
				TestServer.Use(HTTPResponseFactory(func(w http.ResponseWriter) {
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
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
		}

		count := 1
		TestServer.Use(HTTPResponseFactory(func(w http.ResponseWriter) {
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
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X PUT -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X DELETE -H "Content-type:application/json" -d '{"name":"talula"}'`, URLForTestServer("/A")),
			fmt.Sprintf(`%s -X GET`, URLForTestServer("/A")),
		}

		responseBody := "-"
		TestServer.Use(HTTPResponseFactory(func(w http.ResponseWriter) {
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
		for _, method := range HTTPMethodsWithRequestBody {
			It(fmt.Sprintf("in the body for verb %s", method), func() {
				data := "a=1&b=2&c=3"
				list := []string{fmt.Sprintf(`%s -X %s -d %s`, URLForTestServer("/A"), method, data)}
				SutExecute(list)

				predicates := []HTTPRequestPredicate{}
				predicates = append(predicates, RequestWithPath("/A"))
				predicates = append(predicates, RequestWithMethod(method))
				predicates = append(predicates, RequestWithBody(data))
				Expect(TestServer.Find(predicates...)).To(Equal(true))
			})
		}

		It("in the querystring for verb GET", func() {
			method := "GET"
			data := "a=1&b=2&c=3"
			list := []string{fmt.Sprintf(`%s -X %s -d %s"`, URLForTestServer("/A"), method, data)}
			SutExecute(list)

			predicates := []HTTPRequestPredicate{}
			predicates = append(predicates, RequestWithPath("/A"))
			predicates = append(predicates, RequestWithMethod(method))
			predicates = append(predicates, RequestWithQuerystring(data))
			Expect(TestServer.Find(predicates...)).To(Equal(true))
		})
	})

	for _, method := range SupportedHTTPMethods {
		It(fmt.Sprintf("Makes a http %s request with http headers", method), func() {
			applicationJSON := "Content-Type:application/json"
			applicationSoapXML := "Accept:application/soap+xml"
			list := []string{fmt.Sprintf(`%s -X %s -H "%s" -H "%s"`, URLForTestServer("/A"), method, applicationJSON, applicationSoapXML)}
			SutExecute(list)

			predicates := []HTTPRequestPredicate{}
			predicates = append(predicates, RequestWithPath("/A"))
			predicates = append(predicates, RequestWithMethod(method))
			predicates = append(predicates, RequestWithHeader("Content-Type", "application/json"))
			predicates = append(predicates, RequestWithHeader("Accept", "application/soap+xml"))
			Expect(TestServer.Find(predicates...)).To(Equal(true))
		})
	}

	for _, method := range SupportedHTTPMethods {
		It(fmt.Sprintf("Makes a http %s request", method), func() {
			list := []string{fmt.Sprintf(`%s -X %s`, URLForTestServer("/A"), method)}
			SutExecute(list)
			Expect(TestServer.Find(RequestWithPath("/A"), RequestWithMethod(method))).To(Equal(true))
		})
	}

	It("Makes a http get request to each url in a file", func() {
		list := []string{
			URLForTestServer("/A"),
			URLForTestServer("/B"),
			URLForTestServer("/C"),
		}

		SutExecute(list)

		Expect(TestServer.Find(RequestWithPath("/A"))).To(Equal(true))
		Expect(TestServer.Find(RequestWithPath("/B"))).To(Equal(true))
		Expect(TestServer.Find(RequestWithPath("/C"))).To(Equal(true))
	})
})
