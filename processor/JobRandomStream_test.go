package processor_test

import (
	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobRandomStream", func() {

	It("iterates", func() {
		jobs := []core.Job{
			core.Job{Name: "1"},
			core.Job{Name: "2"},
			core.Job{Name: "3"},
		}

		randomStream := CreateJobRandomStream(jobs)
		randomStream.Next()
		randomStream.Next()
		randomStream.Next()
		Expect(randomStream.HasNext()).To(Equal(false))
	})
})
