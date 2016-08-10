package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	It("Does not override set step name", func() {
		expectedName := "fubar"

		job := Job{
			Steps: []Step{},
		}

		step := job.CreateStep()
		step.Name = expectedName

		job = job.AddStep(step)

		Expect(job.Steps[0].Name).To(Equal(expectedName))
	})
})
