package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
