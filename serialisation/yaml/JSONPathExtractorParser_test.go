package yaml_test

import (
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
	. "github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSONPathExtractorParser", func() {
	It("Parses", func() {
		var expectedName = "A"
		var expectedKey = "B"
		var expectedJsonPath = "C"

		var input = map[string]interface{}{
			"name":     expectedName,
			"key":      expectedKey,
			"jsonpath": expectedJsonPath,
			"scope":    core.PlanScope,
		}

		var parser = JSONPathExtractorParser{}

		var result = parser.Parse(input)

		Expect(result).To(BeAssignableToTypeOf(extractors.JSONPathExtractor{}))

		var jsonPathExtractor = result.(extractors.JSONPathExtractor)

		Expect(jsonPathExtractor.Name).To(Equal(expectedName))
		Expect(jsonPathExtractor.Key).To(Equal(expectedKey))
		Expect(jsonPathExtractor.JSONPath).To(Equal(expectedJsonPath))
	})
})
