package converters

import (
	"strings"

	"github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("W3cExtConverter", func() {
	var (
		converter *W3cExtConverter
		plan      *yaml.ExecutionPlan
		err       error
	)

	Describe("when the ", func() {
		BeforeEach(func() {
			const input = `#Fields: date time c-ip cs-username s-computername s-ip cs-method cs-uri-stem cs-uri-query sc-status sc-bytes cs-bytes time-taken cs-version cs(User-Agent) cs(Cookie) cs(Referer)
1996-01-01 10:48:02 195.52.225.44 - WEB1 192.166.0.24 GET /default.htm - 200 1703 279 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) - http://www.webtrends.com/def_f1.htm
1996-01-01 10:48:02 195.52.225.44 - WEB1 192.166.0.24 GET /loganalyzer/info.htm sourceid=chrome-instant&ion=1&espv=2&ie=UTF-8#q=sample%20iis%20log%20files 200 3960 303 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 http://www.webtrends.com/def_f1.htm
1996-01-01 10:48:05 195.52.225.44 - WEB1 192.166.0.24 GET /styles/style1.css - 200 586 249 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 -
1996-01-01 10:48:05 195.52.225.44 - WEB1 192.166.0.24 GET /graphics/atremote/remote.jpg - 200 12367 301 656 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 http://webtrends.sample.com/wt_f2.htm
1996-01-01 10:48:05 195.52.225.44 - WEB1 192.166.0.24 GET /graphics/backg/backg1.gif - 200 448 313 0 HTTP/1.0 Mozilla/4.0+[en]+(WinNT;+I) WEBTRENDS_ID=195.52.225.44-100386000.29188902 http://webtrends.sample.com/loganalyzer/info.htm`
			converter = NewW3cExtConverter("http://mybase.uri", strings.NewReader(input))
			plan, err = converter.Convert()
		})

		It("does not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("builds a plan with one job and many steps", func() {
			Ω(plan.Jobs).Should(HaveLen(1))
			Ω(plan.Jobs[0].Steps).Should(HaveLen(5))
		})

		It("builds a plan with a GET HttpRequest", func() {
			action := plan.Jobs[0].Steps[0].Action
			Ω(action).Should(BeAssignableToTypeOf(yaml.Action{}))
			Ω(action["type"]).Should(Equal("HttpRequest"))
			Ω(action["method"]).Should(Equal("GET"))
		})

		It("adds an ExactAssertion for the HTTP status", func() {
			assertion := plan.Jobs[0].Steps[0].Assertions[0]
			Ω(assertion).Should(BeAssignableToTypeOf(yaml.Assertion{}))
			Ω(assertion["type"]).Should(Equal("ExactAssertion"))
			Ω(assertion["key"]).Should(Equal("response:status"))
		})
	})
})
