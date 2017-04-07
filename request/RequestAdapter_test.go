package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/guzzlerio/corcel/global"

	. "github.com/smartystreets/goconvey/convey"
)

func URLForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", global.TestPort, path)
}

func TestRequestAdapter(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("RequestAdapter", t, func() {
		var (
			userAgent = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 5_1_1 like Mac OS X; en) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3"
			url       string
			line      string
			adapter   Adapter
			req       *http.Request
			err       error
		)

		func() {
			url = "http://localhost:8000/A"
			line = url
			line += " -X POST"
			line += ` -H "Content-type: application/json"`
			line += fmt.Sprintf(` -A "%s"`, userAgent)
			adapter = NewRequestAdapter()
			reqFunc := adapter.Create(line)
			req, err = reqFunc()
			So(err, ShouldBeNil)
		}()

		Convey("Parses URL", func() {
			So(req.URL.String(), ShouldEqual, url)
		})

		Convey("Parses Method", func() {
			So(req.Method, ShouldEqual, "POST")
		})

		Convey("Parses Header", func() {
			So(req.Header.Get("Content-type"), ShouldEqual, "application/json")
		})

		Convey("Unhappy path", func() {
			Convey("Does something", func() {
				adapter = NewRequestAdapter()
				reqFunc := adapter.Create("-Something")
				_, err := reqFunc()
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Parses Body", func() {
			Convey("For GET request is inside the querystring", func() {
				data := "a=1&b=2"
				line = url
				line += " -X GET"
				line += fmt.Sprintf(` -d "%s"`, data)
				adapter = NewRequestAdapter()
				reqFunc := adapter.Create(line)
				req, err = reqFunc()
				So(err, ShouldBeNil)
				So(req.URL.RawQuery, ShouldEqual, data)
			})

			for _, method := range global.HTTPMethodsWithRequestBody {
				Convey(fmt.Sprintf("For %s request is in the actual request body", method), func() {
					data := "a=1&b=2"
					line = url
					line += fmt.Sprintf(" -X %s", method)
					line += fmt.Sprintf(` -d "%s"`, data)
					adapter = NewRequestAdapter()
					reqFunc := adapter.Create(line)
					req, err = reqFunc()
					So(err, ShouldBeNil)
					body, bodyErr := ioutil.ReadAll(req.Body)
					check(bodyErr)
					So(string(body), ShouldEqual, data)
				})

				Convey(fmt.Sprintf("For %s requests specifying an input file", method), func() {
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
					So(err, ShouldBeNil)
					body, bodyErr := ioutil.ReadAll(req.Body)
					check(bodyErr)
					So(string(body), ShouldEqual, data)

				})
			}

		})

		Convey("Parses URLs with leading whitespace", func() {
			line = "      " + url
			adapter = NewRequestAdapter()
			reqFunc := adapter.Create(line)
			req, err = reqFunc()
			So(req.URL.String(), ShouldEqual, url)
		})

		Convey("Parses UserAgent", func() {
			So(req.UserAgent(), ShouldEqual, userAgent)
		})
	})
}
