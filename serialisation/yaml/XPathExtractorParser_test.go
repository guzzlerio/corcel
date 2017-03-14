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

		var result, err = parser.Parse(input)

		Expect(err).To(BeNil())

		Expect(result).To(BeAssignableToTypeOf(extractors.XPathExtractor{}))

		var xpathExtractor = result.(extractors.XPathExtractor)
		Expect(xpathExtractor.Key).To(Equal(expectedKey))
		Expect(xpathExtractor.Name).To(Equal(expectedName))
		Expect(xpathExtractor.XPath).To(Equal(expectedXpath))
		Expect(xpathExtractor.Scope).To(Equal(core.PlanScope))
	})

	It("Fails to parse with empty name", func() {
		var input = map[string]interface{}{
			"key":   "key",
			"xpath": "path",
			"scope": core.PlanScope,
		}
		var parser = XPathExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})

	It("Fails to parse with empty key", func() {
		var input = map[string]interface{}{
			"name":  "name",
			"xpath": "path",
			"scope": core.PlanScope,
		}
		var parser = XPathExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})

	It("Fails to parse with empty xpath", func() {
		var input = map[string]interface{}{
			"key":   "key",
			"name":  "name",
			"scope": core.PlanScope,
		}
		var parser = XPathExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})
})
