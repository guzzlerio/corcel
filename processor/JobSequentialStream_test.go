package processor

import (
	"github.com/guzzlerio/corcel/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobSequentialStream", func() {

	It("iterates", func() {
		jobs := []core.Job{
			core.Job{Name: "1"},
			core.Job{Name: "2"},
			core.Job{Name: "3"},
		}

		sequentialStream := CreateJobSequentialStream(jobs)
		Expect(sequentialStream.Next()).To(Equal(jobs[0]))
		Expect(sequentialStream.Next()).To(Equal(jobs[1]))
		Expect(sequentialStream.Next()).To(Equal(jobs[2]))
		Expect(sequentialStream.HasNext()).To(Equal(false))
	})

})
