package processor_test

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/processor"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJobRandomStream(t *testing.T) {
	BeforeTest()
	defer AfterTest()
	Convey("JobRandomStream", t, func() {

		Convey("iterates", func() {
			jobs := []core.Job{
				core.Job{Name: "1"},
				core.Job{Name: "2"},
				core.Job{Name: "3"},
			}

			randomStream := CreateJobRandomStream(jobs)
			randomStream.Next()
			randomStream.Next()
			randomStream.Next()
			So(randomStream.HasNext(), ShouldEqual, false)
		})
	})
}
