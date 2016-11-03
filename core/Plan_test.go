package core_test

import (
	. "github.com/guzzlerio/corcel/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {
	It("Does not override set job name", func() {
		expectedName := "talula"
		plan := Plan{
			Jobs: []Job{},
		}

		job := plan.CreateJob()
		job.Name = expectedName

		plan = plan.AddJob(job)
		Expect(plan.Jobs[0].Name).To(Equal(expectedName))
	})
})
