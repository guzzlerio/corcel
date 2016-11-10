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
			const input = `#Version: 1.0
#Date: 12-Jan-1996 00:00:00
#Fields: time cs-method cs-uri
00:34:23 GET /foo/bar.html`
			converter = NewW3cExtConverter("http://mybase.uri", strings.NewReader(input))
			plan, err = converter.Convert()
		})

		It("Does not override set step name", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("builds a plan with a GET HttpRequest", func() {
			Ω(plan.Jobs[0].Steps[0].Action).Should(BeAssignableToTypeOf(yaml.Action{}))
		})
	})
})
