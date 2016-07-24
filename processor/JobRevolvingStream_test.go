package processor_test

import (
	"ci.guzzler.io/guzzler/corcel/core"
	. "ci.guzzler.io/guzzler/corcel/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobRevolvingStream", func() {

	It("iterates", func() {
		jobs := []core.Job{
			core.Job{Name: "1"},
			core.Job{Name: "2"},
			core.Job{Name: "3"},
		}

		sequentialStream := CreateJobSequentialStream(jobs)
		revolvingStream := CreateJobRevolvingStream(sequentialStream)
		Expect(revolvingStream.Next()).To(Equal(jobs[0]))
		Expect(revolvingStream.Next()).To(Equal(jobs[1]))
		Expect(revolvingStream.Next()).To(Equal(jobs[2]))
		Expect(revolvingStream.Next()).To(Equal(jobs[0]))
	})
})
