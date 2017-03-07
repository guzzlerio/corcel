package yaml

import (
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("XPathExtractorParser", func() {
	It("Parses", func() {
		var expectedName = "A"
		var expectedKey = "B"
		var expectedXpath = "C"

		var input = map[string]interface{}{
			"name":  expectedName,
			"key":   expectedKey,
			"xpath": expectedXpath,
			"scope": core.PlanScope,
		}

		var parser = XPathExtractorParser{}

		var result = parser.Parse(input)

		Expect(result).To(BeAssignableToTypeOf(extractors.XPathExtractor{}))

		var xpathExtractor = result.(extractors.XPathExtractor)
		Expect(xpathExtractor.Key).To(Equal(expectedKey))
		Expect(xpathExtractor.Name).To(Equal(expectedName))
		Expect(xpathExtractor.XPath).To(Equal(expectedXpath))
		Expect(xpathExtractor.Scope).To(Equal(core.PlanScope))
	})
})
