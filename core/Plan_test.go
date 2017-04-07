package core_test

import (
	"testing"

	. "github.com/guzzlerio/corcel/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPlan(t *testing.T) {
	Convey("Plan", t, func() {
		Convey("Does not override set job name", func() {
			expectedName := "talula"
			plan := Plan{
				Jobs: []Job{},
			}

			job := plan.CreateJob()
			job.Name = expectedName

			plan = plan.AddJob(job)
			So(plan.Jobs[0].Name, ShouldEqual, expectedName)
		})
	})
}
