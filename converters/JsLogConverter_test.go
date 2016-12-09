package converters

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("JsLogConverter", func() {
	Describe("iisParser", func() {
		const input = `#Fields: date time c-ip cs-username s-computername s-ip cs-method cs-uri-stem cs-uri-query sc-status sc-bytes cs-bytes time-taken cs-version cs(User-Agent) cs(Cookie) cs(Referer)
1996-01-01 10:48:02 195.52.225.44 - WEB1 192.166.0.24 GET /default.htm - 200 1703 279 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) - http://www.webtrends.com/def_f1.htm
1996-01-01 10:48:02 195.52.225.44 - WEB1 192.166.0.24 GET /loganalyzer/info.htm sourceid=chrome-instant&ion=1&espv=2&ie=UTF-8#q=sample%20iis%20log%20files 200 3960 303 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 http://www.webtrends.com/def_f1.htm
1996-01-01 10:48:05 195.52.225.44 - WEB1 192.166.0.24 GET /styles/style1.css - 200 586 249 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 -
1996-01-01 10:48:05 195.52.225.44 - WEB1 192.166.0.24 GET /graphics/atremote/remote.jpg - 200 12367 301 656 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 http://webtrends.sample.com/wt_f2.htm
1996-01-01 10:48:05 195.52.225.44 - WEB1 192.166.0.24 GET /graphics/backg/backg1.gif - 200 448 313 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 http://webtrends.sample.com/loganalyzer/info.htm`
		var (
			converter *JsLogConverter
			plan      *yaml.ExecutionPlan
			err       error
			parser    string
			baseUrl   *url.URL
		)

		BeforeSuite(func() {
			buf, _ := ioutil.ReadFile("./parsers/iisParser.js")
			parser = string(buf)
			baseUrl, _ = url.Parse("http://blah.com")
		})

		Describe("when the js parser is valid javascript", func() {
			BeforeEach(func() {
				converter = NewJsLogConverter(parser, baseUrl, strings.NewReader(input))
				plan, err = converter.Convert()
				// WriteOutputYAML(plan)
			})

			It("does not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("builds a plan", func() {
				Ω(plan).Should(BeAssignableToTypeOf(&yaml.ExecutionPlan{}))
			})

			Describe("it builds the jobs", func() {
				It("and builds one job", func() {
					Ω(plan.Jobs).Should(HaveLen(1))
				})
			})

			Describe("the Steps", func() {
				It("are all added", func() {
					Ω(plan.Jobs[0].Steps).Should(HaveLen(5))
				})

				It("are built with a GET HttpRequest", func() {
					action := plan.Jobs[0].Steps[1].Action
					Ω(action).Should(BeAssignableToTypeOf(yaml.Action{}))
					Ω(action["type"]).Should(Equal("HttpRequest"))
					Ω(action["method"]).Should(Equal("GET"))
					Ω(action["url"]).Should(Equal(baseUrl.String() + "/loganalyzer/info.htm?sourceid=chrome-instant&ion=1&espv=2&ie=UTF-8#q=sample%20iis%20log%20files"))
				})

				It("add an ExactAssertion for the HTTP status", func() {
					assertion := plan.Jobs[0].Steps[0].Assertions[0]
					Ω(assertion).Should(BeAssignableToTypeOf(yaml.Assertion{}))
					Ω(assertion["type"]).Should(Equal("ExactAssertion"))
					Ω(assertion["key"]).Should(Equal("response:status"))
					Ω(assertion["expected"]).Should(Equal(200))
				})
			})

			Describe("and the input log file contains POST requests", func() {
				Describe("but does not contain the payload", func() {
					BeforeEach(func() {
						postLine := "1996-01-01 10:48:02 195.52.225.44 - WEB1 192.166.0.24 POST /default.htm - 200 1703 279 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) - http://www.webtrends.com/def_f1.htm"
						talula := input + "\n" + postLine
						converter = NewJsLogConverter(parser, baseUrl, strings.NewReader(talula))
						plan, err = converter.Convert()
					})

					It("it does not add a Step for that line", func() {
						for _, step := range plan.Jobs[0].Steps {
							Ω(step.Action).ShouldNot(HaveKeyWithValue("method", "POST"))
						}
					})
				})
			})

			Describe("and the input log file contains a new field list", func() {
				BeforeEach(func() {
					postLine := `#Fields: date time cs-method cs-uri-stem cs-uri-query sc-status sc-bytes cs-bytes time-taken cs-version cs(User-Agent) cs(Cookie) cs(Referer)
1996-01-01 10:48:02 GET /default.htm - 200 1703 279 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) - http://www.webtrends.com/def_f1.htm`
					talula := input + "\n" + postLine
					converter = NewJsLogConverter(parser, baseUrl, strings.NewReader(talula))
					plan, err = converter.Convert()
				})

				It("does not error", func() {
					Ω(err).ShouldNot(HaveOccurred())
				})
				It("are all added", func() {
					Ω(plan.Jobs[0].Steps).Should(HaveLen(6))
				})
			})

			Describe("but the input log file contains invalid data", func() {
				It("Blows up!", func() {
					talula := `************`
					converter = NewJsLogConverter(parser, baseUrl, strings.NewReader(talula))
					plan, err = converter.Convert()
				})
			})
			Describe("but the input log file does not contain sufficient fields", func() {})
		})

		Describe("when the js is invalid", func() {
			It("panics", func() {
				defer func() {
					err := recover()
					Ω(err).Should(HaveOccurred())
				}()
				NewJsLogConverter("not valid javascript", baseUrl, strings.NewReader(input))
			})
		})

		Describe("when the baseUrl is not supplied", func() {
			It("panics", func() {
				defer func() {
					err := recover()
					Ω(err).Should(HaveOccurred())
				}()
				NewJsLogConverter(parser, nil, strings.NewReader(input))
			})
		})

		Describe("when the js errors", func() {
			It("errors", func() {
				parser = `function parseLine(line, fields) { throw new Error('BOOOOOM'); }`
				converter = NewJsLogConverter(parser, baseUrl, strings.NewReader(input))
				plan, err := converter.Convert()
				Ω(err).Should(HaveOccurred())
				Ω(plan).Should(BeNil())
			})
		})
	})
})

func WriteOutputYAML(plan *yaml.ExecutionPlan) {
	planBuilder := yaml.NewPlanBuilder()
	if file, err := planBuilder.Write(plan); err == nil {
		defer func() {
			if fileErr := os.Remove(file.Name()); fileErr != nil {
				panic(fileErr)
			}
		}()
		dat, _ := ioutil.ReadFile(file.Name())

		fmt.Println(string(dat))
	}
}
