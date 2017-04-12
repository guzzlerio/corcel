package processor_test

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/processor"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJobIterationStream(t *testing.T) {
	BeforeTest()
	defer AfterTest()

	Convey("JobIterationStream", t, func() {
		Convey("iterates", func() {
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
				So(iterationStream.Next(), ShouldResemble, jobs[i%len(jobs)])
			}
			So(iterationStream.HasNext(), ShouldEqual, false)
		})
	})
}
