package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	//SupportedHTTPMethods ...
	SupportedHTTPMethods = []string{"GET", "POST", "PUT", "DELETE"}
	//HTTPMethodsWithRequestBody ...
	HTTPMethodsWithRequestBody = []string{"POST", "PUT", "DELETE"}
	//TestServer ...
	TestServer *RequestRecordingServer
	//TestPort ...
	TestPort = 8000
	//ResponseCodes400 ...
	ResponseCodes400 = []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418}
	//ResponseCodes500 ...
	ResponseCodes500 = []int{500, 501, 502, 503, 504, 505}
	//WaitTimeTests ...
	WaitTimeTests = []string{"1ms", "2ms", "4ms", "8ms", "16ms", "32ms", "64ms", "128ms"}
	//NumberOfWorkersToTest ...
	NumberOfWorkersToTest = []int{1, 2, 4, 8, 16, 32, 64, 128, 256}
)

func URLForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TestPort, path)
}

var _ = Describe("RequestAdapter", func() {
	var (
		userAgent = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 5_1_1 like Mac OS X; en) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3"
		url       string
		line      string
		adapter   RequestAdapter
		req       *http.Request
		err       error
	)

	BeforeEach(func() {
		url = "http://localhost:8000/A"
		line = url
		line += " -X POST"
		line += ` -H "Content-type: application/json"`
		line += fmt.Sprintf(` -A "%s"`, userAgent)
		adapter = NewRequestAdapter()
		reqFunc := adapter.Create(line)
		req, err = reqFunc()
		Expect(err).To(BeNil())
	})

	It("Parses URL", func() {
		Expect(req.URL.String()).To(Equal(url))
	})

	It("Parses Method", func() {
		Expect(req.Method).To(Equal("POST"))
	})

	It("Parses Header", func() {
		Expect(req.Header.Get("Content-type")).To(Equal("application/json"))
	})

	Describe("Parses Body", func() {
		It("For GET request is inside the querystring", func() {
			data := "a=1&b=2"
			line = url
			line += " -X GET"
			line += fmt.Sprintf(` -d "%s"`, data)
			adapter = NewRequestAdapter()
			reqFunc := adapter.Create(line)
			req, err = reqFunc()
			Expect(err).To(BeNil())
			Expect(req.URL.RawQuery).To(Equal(data))
		})

		for _, method := range HTTPMethodsWithRequestBody {
			It(fmt.Sprintf("For %s request is in the actual request body", method), func() {
				data := "a=1&b=2"
				line = url
				line += fmt.Sprintf(" -X %s", method)
				line += fmt.Sprintf(` -d "%s"`, data)
				adapter = NewRequestAdapter()
				reqFunc := adapter.Create(line)
				req, err = reqFunc()
				Expect(err).To(BeNil())
				body, bodyErr := ioutil.ReadAll(req.Body)
				check(bodyErr)
				Expect(string(body)).To(Equal(data))
			})

			It(fmt.Sprintf("For %s requests specifying an input file", method), func() {
				data := "a=1&b=2"
				loadRequestBodyFromFile = func(filename string) *bytes.Buffer {
					body := bytes.NewBuffer([]byte(data))
					return body
				}

				line = url
				line += fmt.Sprintf(" -X %s", method)
				line += " -d @./file"
				adapter = NewRequestAdapter()
				reqFunc := adapter.Create(line)
				req, err = reqFunc()
				Expect(err).To(BeNil())
				body, bodyErr := ioutil.ReadAll(req.Body)
				check(bodyErr)
				Expect(string(body)).To(Equal(data))

			})
		}

	})

	It("Parses URLs with leading whitespace", func() {
		line = "      " + url
		adapter = NewRequestAdapter()
		reqFunc := adapter.Create(line)
		req, err = reqFunc()
		Expect(req.URL.String()).To(Equal(url))
	})

	It("Parses UserAgent", func() {
		Expect(req.UserAgent()).To(Equal(userAgent))
	})
})
