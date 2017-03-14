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

		var result, err = parser.Parse(input)

		Expect(err).To(BeNil())

		Expect(result).To(BeAssignableToTypeOf(extractors.JSONPathExtractor{}))

		var jsonPathExtractor = result.(extractors.JSONPathExtractor)

		Expect(jsonPathExtractor.Name).To(Equal(expectedName))
		Expect(jsonPathExtractor.Key).To(Equal(expectedKey))
		Expect(jsonPathExtractor.JSONPath).To(Equal(expectedJsonPath))
	})

	It("Fails to parse when name not set", func() {

		var input = map[string]interface{}{
			"key":      "key",
			"jsonpath": "path",
			"scope":    core.PlanScope,
		}

		var parser = JSONPathExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})

	It("Fails to parse when key not set", func() {

		var input = map[string]interface{}{
			"name":     "name",
			"jsonpath": "path",
			"scope":    core.PlanScope,
		}

		var parser = JSONPathExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})

	It("Fails to parse when jsonpath not set", func() {

		var input = map[string]interface{}{
			"name":  "name",
			"key":   "key",
			"scope": core.PlanScope,
		}

		var parser = JSONPathExtractorParser{}

		var _, err = parser.Parse(input)

		Expect(err).ToNot(BeNil())
	})
})
