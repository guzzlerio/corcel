package yaml_test

import (
	"github.com/guzzlerio/corcel/extractors"
	. "github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyValueExtractorParser", func() {
	It("Parses", func() {
		var expectedName = "A"
		var expectedKey = "B"

		var input = map[string]interface{}{
			"name": expectedName,
			"key":  expectedKey,
		}

		var parser = KeyValueExtractorParser{}

		var result = parser.Parse(input)

		Expect(result).To(BeAssignableToTypeOf(extractors.KeyValueExtractor{}))

		var keyValueExtractor = result.(extractors.KeyValueExtractor)

		Expect(keyValueExtractor.Name).To(Equal(expectedName))
		Expect(keyValueExtractor.Key).To(Equal(expectedKey))
	})
})
