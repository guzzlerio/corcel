package extractors

import (
	"github.com/guzzlerio/corcel/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyValueExtractor", func() {

	It("Succeeds when key is present", func() {
		var extractor = KeyValueExtractor{
			Key:   "key",
			Name:  "target",
			Scope: core.StepScope,
		}

		var executionResult = core.ExecutionResult{
			"key": "talula",
		}

		var extractionResult = extractor.Extract(executionResult)

		Expect(extractionResult["target"]).To(Equal("talula"))
		Expect(extractionResult["scope"]).To(Equal(core.StepScope))
	})

	It("Extraction result does not contain the name when the key does not exist inside the execution result", func() {
		var extractor = KeyValueExtractor{
			Key:   "key",
			Name:  "target",
			Scope: core.StepScope,
		}

		var executionResult = core.ExecutionResult{
			"lock": "talula",
		}

		var extractionResult = extractor.Extract(executionResult)

		Expect(extractionResult["target"]).To(BeNil())
	})

})
