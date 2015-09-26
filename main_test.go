package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	SUPPORTED_HTTP_METHODS         = []string{"GET", "POST", "PUT", "DELETE"}
	HTTP_METHODS_WITH_REQUEST_BODY = []string{"POST", "PUT", "DELETE"}
	server                         *RequestRecordingServer
	port                           int
)

func TestMain(m *testing.M) {
	configureLogging()
	port = 8000
	server = CreateRequestRecordingServer(port)
	server.Start()
	os.Exit(m.Run())
	server.Stop()
}

func CreateList(lines []string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		file.WriteString(fmt.Sprintf("%s\n", line))
	}
	file.Sync()
	return file
}

func PathExists(value string) bool {
	path, pathErr := filepath.Abs(value)
	if pathErr != nil{
		panic(pathErr)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

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
		server.Clear()
	})

	FIt("Generate statistics of data from the execution", func() {
		list := []string{
			`http://127.0.0.1:8000/A -X POST -H "Content-type:application/json" -d '{"name":"talula"}'`,
			`http://127.0.0.1:8000/A -X PUT -H "Content-type:application/json" -d '{"name":"talula"}'`,
			`http://127.0.0.1:8000/A -X DELETE -H "Content-type:application/json" -d '{"name":"talula"}'`,
			`http://127.0.0.1:8000/A -X GET`,
		}

		file := CreateList(list)
		defer os.Remove(file.Name())
		cmd := exec.Command(exePath, "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		Expect(err).To(BeNil())

		Expect(PathExists("./output.json")).To(Equal(true))

	})

	Describe("Support sending data with http request", func() {
		for _, method := range HTTP_METHODS_WITH_REQUEST_BODY {
			It(fmt.Sprintf("in the body for verb %s", method), func() {
				data := "a=1&b=2&c=3"
				list := []string{fmt.Sprintf(`http://127.0.0.1:8000/A -X %s -d %s`, method, data)}
				file := CreateList(list)
				defer os.Remove(file.Name())
				cmd := exec.Command(exePath, "-f", file.Name())
				output, err := cmd.CombinedOutput()
				fmt.Println(string(output))
				Expect(err).To(BeNil())

				predicates := []HttpRequestPredicate{}
				predicates = append(predicates, RequestWithPath("/A"))
				predicates = append(predicates, RequestWithMethod(method))
				predicates = append(predicates, RequestWithBody(data))
				Expect(server.Find(predicates...)).To(Equal(true))
			})
		}

		It("in the querystring for verb GET", func() {
			method := "GET"
			data := "a=1&b=2&c=3"
			list := []string{fmt.Sprintf(`http://127.0.0.1:8000/A -X %s -d %s"`, method, data)}
			file := CreateList(list)
			defer os.Remove(file.Name())
			cmd := exec.Command(exePath, "-f", file.Name())
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			Expect(err).To(BeNil())

			predicates := []HttpRequestPredicate{}
			predicates = append(predicates, RequestWithPath("/A"))
			predicates = append(predicates, RequestWithMethod(method))
			predicates = append(predicates, RequestWithQuerystring(data))
			Expect(server.Find(predicates...)).To(Equal(true))
		})
	})

	for _, method := range SUPPORTED_HTTP_METHODS {
		It(fmt.Sprintf("Makes a http %s request with http headers", method), func() {
			applicationJson := "Content-Type:application/json"
			applicationSoapXml := "Accept:application/soap+xml"
			list := []string{fmt.Sprintf(`http://127.0.0.1:8000/A -X %s -H "%s" -H "%s"`, method, applicationJson, applicationSoapXml)}
			file := CreateList(list)
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
			Expect(server.Find(predicates...)).To(Equal(true))
		})
	}

	for _, method := range SUPPORTED_HTTP_METHODS {
		It(fmt.Sprintf("Makes a http %s request", method), func() {
			list := []string{fmt.Sprintf(`http://127.0.0.1:8000/A -X %s`, method)}
			file := CreateList(list)
			defer os.Remove(file.Name())
			cmd := exec.Command(exePath, "-f", file.Name())
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			Expect(err).To(BeNil())
			Expect(server.Find(RequestWithPath("/A"), RequestWithMethod(method))).To(Equal(true))
		})
	}

	It("Makes a http get request to each url in a file", func() {
		list := []string{"http://127.0.0.1:8000/A",
			"http://127.0.0.1:8000/B",
			"http://127.0.0.1:8000/C"}
		file := CreateList(list)
		defer os.Remove(file.Name())

		cmd := exec.Command(exePath, "-f", file.Name())
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		Expect(err).To(BeNil())
		Expect(server.Find(RequestWithPath("/A"))).To(Equal(true))
		Expect(server.Find(RequestWithPath("/B"))).To(Equal(true))
		Expect(server.Find(RequestWithPath("/C"))).To(Equal(true))
	})
})
