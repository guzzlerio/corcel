package processor

import (
	"testing"

	"github.com/guzzlerio/corcel/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJobSequentialStream(t *testing.T) {
	BeforeTest()
	defer AfterTest()
	Convey("JobSequentialStream", t, func() {

		Convey("iterates", func() {
			jobs := []core.Job{
				core.Job{Name: "1"},
				core.Job{Name: "2"},
				core.Job{Name: "3"},
			}

			sequentialStream := CreateJobSequentialStream(jobs)
			So(sequentialStream.Next(), ShouldResemble, jobs[0])
			So(sequentialStream.Next(), ShouldResemble, jobs[1])
			So(sequentialStream.Next(), ShouldResemble, jobs[2])
			So(sequentialStream.HasNext(), ShouldEqual, false)
		})

	})
}
