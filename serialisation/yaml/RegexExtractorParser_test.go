package yaml_test

import (
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
	. "github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegexExtractorParser", func() {
	It("Parses", func() {
		var expectedName = "A"
		var expectedKey = "B"
		var expectedMatch = "C"

		var input = map[string]interface{}{
			"name":  expectedName,
			"key":   expectedKey,
			"match": expectedMatch,
			"scope": core.PlanScope,
		}

		var parser = RegexExtractorParser{}

		var result, err = parser.Parse(input)

		Expect(err).To(BeNil())

		Expect(result).To(BeAssignableToTypeOf(extractors.RegexExtractor{}))

		var jsonPathExtractor = result.(extractors.RegexExtractor)

		Expect(jsonPathExtractor.Name).To(Equal(expectedName))
		Expect(jsonPathExtractor.Key).To(Equal(expectedKey))
		Expect(jsonPathExtractor.Match).To(Equal(expectedMatch))
	})

	It("Fails to parse without name", func() {

		var input = map[string]interface{}{
			"key":   "key",
			"match": "match",
			"scope": core.PlanScope,
		}

		var parser = RegexExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})

	It("Fails to parse without key", func() {

		var input = map[string]interface{}{
			"name":  "name",
			"match": "match",
			"scope": core.PlanScope,
		}

		var parser = RegexExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})

	It("Fails to parse without match", func() {

		var input = map[string]interface{}{
			"name":  "name",
			"key":   "key",
			"scope": core.PlanScope,
		}

		var parser = RegexExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})
})
