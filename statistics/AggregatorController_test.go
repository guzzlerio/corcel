package statistics_test

import (
	"testing"
	"time"

	. "github.com/guzzlerio/corcel/statistics"
	metrics "github.com/rcrowley/go-metrics"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregatorController(t *testing.T) {
	Convey("AggregatorController", t, func() {

		Convey("Can take a snapshot after calling stop", func() {
			var controller = CreateAggregatorController(metrics.DefaultRegistry)

			var snapshot AggregatorSnapShot

			controller.Start()
			controller.Stop()
			snapshot = controller.Snapshot()

			So(len(snapshot.Times), ShouldEqual, 2)
		})

		Convey("Creates a snapshot after a time period", func() {
			var controller = CreateAggregatorController(metrics.DefaultRegistry)

			var snapshot AggregatorSnapShot

			controller.Start()
			time.Sleep(2 * time.Second)
			controller.Stop()
			snapshot = controller.Snapshot()

			So(len(snapshot.Times), ShouldEqual, 3)
		})
	})
}
