package processor_test

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/processor"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJobRevolvingStream(t *testing.T) {
	BeforeTest()
	defer AfterTest()
	Convey("JobRevolvingStream", t, func() {

		Convey("iterates", func() {
			jobs := []core.Job{
				core.Job{Name: "1"},
				core.Job{Name: "2"},
				core.Job{Name: "3"},
			}

			sequentialStream := CreateJobSequentialStream(jobs)
			revolvingStream := CreateJobRevolvingStream(sequentialStream)
			So(revolvingStream.Next(), ShouldResemble, jobs[0])
			So(revolvingStream.Next(), ShouldResemble, jobs[1])
			So(revolvingStream.Next(), ShouldResemble, jobs[2])
			So(revolvingStream.Next(), ShouldResemble, jobs[0])
		})
	})
}
