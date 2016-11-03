package processor_test

import (
	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobIterationStream", func() {
	It("iterates", func() {
		jobs := []core.Job{
			core.Job{Name: "1"},
			core.Job{Name: "2"},
			core.Job{Name: "3"},
		}

		iterations := 5

		sequentialStream := CreateJobSequentialStream(jobs)
		revolvingStream := CreateJobRevolvingStream(sequentialStream)
		iterationStream := CreateJobIterationStream(*revolvingStream, len(jobs), iterations)

		for i := 0; i < iterations*len(jobs); i++ {
			Expect(iterationStream.Next()).To(Equal(jobs[i%len(jobs)]))
		}
		Expect(iterationStream.HasNext()).To(Equal(false))
	})
})
