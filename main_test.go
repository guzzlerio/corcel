package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	SUPPORTED_HTTP_METHODS         = []string{"GET", "POST", "PUT", "DELETE"}
	HTTP_METHODS_WITH_REQUEST_BODY = []string{"POST", "PUT", "DELETE"}
	TestServer                     *RequestRecordingServer
	TEST_PORT                      = 8000
)

func UrlForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TEST_PORT, path)
}

var _ = BeforeSuite(func() {
	configureLogging()
	TestServer = CreateRequestRecordingServer(TEST_PORT)
	TestServer.Start()
})

var _ = AfterSuite(func() {
	TestServer.Stop()
})

var _ = Describe("Main", func() {

	var (
		exePath string
		err     error
	)

	BeforeEach(func() {
		exePath, err = filepath.Abs("./code-named-something")
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	It("Generate statistics of data from the execution", func() {
		list := []string{
			fmt.Sprintf(`%s -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X PUT -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X DELETE -H "Content-type:application/json" -d '{"name":"talula"}'`, UrlForTestServer("/A")),
			fmt.Sprintf(`%s -X GET`, UrlForTestServer("/A")),
		}

		file := CreateFileFromLines(list)
		defer os.Remove(file.Name())
		cmd := exec.Command(exePath, "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		Expect(err).To(BeNil())

		Expect(PathExists("./output.yml")).To(Equal(true))

		var executionOutput ExecutionOutput

		UnmarshalYamlFromFile("./output.yml", &executionOutput)

		Expect(executionOutput.Summary.TotalBytesSent).To(Equal(uint64(1)))
	})

	Describe("Support sending data with http request", func() {
		for _, method := range HTTP_METHODS_WITH_REQUEST_BODY {
			It(fmt.Sprintf("in the body for verb %s", method), func() {
				data := "a=1&b=2&c=3"
				list := []string{fmt.Sprintf(`%s -X %s -d %s`, UrlForTestServer("/A"), method, data)}
				file := CreateFileFromLines(list)
				defer os.Remove(file.Name())
				cmd := exec.Command(exePath, "-f", file.Name())
				output, err := cmd.CombinedOutput()
				fmt.Println(string(output))
				Expect(err).To(BeNil())

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
			file := CreateFileFromLines(list)
			defer os.Remove(file.Name())
			cmd := exec.Command(exePath, "-f", file.Name())
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			Expect(err).To(BeNil())

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
			file := CreateFileFromLines(list)
			defer os.Remove(file.Name())
			cmd := exec.Command(exePath, "-f", file.Name())
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			Expect(err).To(BeNil())

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
			file := CreateFileFromLines(list)
			defer os.Remove(file.Name())
			cmd := exec.Command(exePath, "-f", file.Name())
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			Expect(err).To(BeNil())
			Expect(TestServer.Find(RequestWithPath("/A"), RequestWithMethod(method))).To(Equal(true))
		})
	}

	It("Makes a http get request to each url in a file", func() {
		list := []string{
			UrlForTestServer("/A"),
			UrlForTestServer("/B"),
			UrlForTestServer("/C"),
		}
		file := CreateFileFromLines(list)
		defer os.Remove(file.Name())

		cmd := exec.Command(exePath, "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		Expect(err).To(BeNil())
		Expect(TestServer.Find(RequestWithPath("/A"))).To(Equal(true))
		Expect(TestServer.Find(RequestWithPath("/B"))).To(Equal(true))
		Expect(TestServer.Find(RequestWithPath("/C"))).To(Equal(true))
	})
})
