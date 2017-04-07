package statistics

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregatorSummary(t *testing.T) {
	Convey("AggregatorSummary", t, func() {

		var (
			snapshot *AggregatorSnapShot
		)

		func() {
			snapshot = NewAggregatorSnapShot()
		}()

		Convey("handles an empty snapshot", func() {
			So(snapshot.CreateSummary(), ShouldNotBeNil)
		})

		Convey("setting the RunningTime", func() {
			Convey("it calculates the correct duration", func() {
				start := int64(1483228800000)
				end := int64(1483229010000)
				duration := time.Unix(0, end).Sub(time.Unix(0, start))
				snapshot.updateTime(start)
				snapshot.updateTime(end)
				summary := snapshot.CreateSummary()
				So(summary.RunningTime, ShouldEqual, duration)
			})

			Convey("does not permit a negative result", func() {
				//Invert the times from previous test
				end := int64(1483228800000)
				start := int64(1483229010000)
				snapshot.updateTime(start)
				snapshot.updateTime(end)
				summary := snapshot.CreateSummary()
				So(summary.RunningTime, ShouldEqual, time.Duration(0))
			})
		})
	})
}
