package extractors

import (
	"github.com/guzzlerio/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSONPathExtractor", func() {

	It("extracts a single key", func() {
		var extractor = JSONPathExtractor{
			Name:     "something",
			Key:      "targetKey",
			JSONPath: "$.aKey",
		}

		var result = core.ExecutionResult{
			"targetKey": `{"aKey":32}`,
		}

		var extractionResult = extractor.Extract(result)

		Expect(extractionResult["something"]).To(Equal(float64(32)))
	})

	It("invalid json path", func() {
		var extractor = JSONPathExtractor{
			Name:     "something",
			Key:      "targetKey",
			JSONPath: "talula",
		}

		var result = core.ExecutionResult{
			"targetKey": `{"aKey":32}`,
		}

		var extractionResult = extractor.Extract(result)

		Expect(extractionResult["something"]).To(Equal(ErrInvalidJsonPath))
	})
})
