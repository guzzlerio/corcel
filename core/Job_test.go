package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJob(t *testing.T) {
	Convey("Job", t, func() {
		Convey("Does not override set step name", func() {
			expectedName := "fubar"

			job := Job{
				Steps: []Step{},
			}

			step := job.CreateStep()
			step.Name = expectedName

			job = job.AddStep(step)

			So(job.Steps[0].Name, ShouldEqual, expectedName)
		})
	})
}
