package statistics_test

import (
	"time"

	. "github.com/guzzlerio/corcel/statistics"
	metrics "github.com/rcrowley/go-metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AggregatorController", func() {

	It("Can take a snapshot after calling stop", func() {
		var controller = CreateAggregatorController(metrics.DefaultRegistry)

		var snapshot AggregatorSnapShot

		controller.Start()
		controller.Stop()
		snapshot = controller.Snapshot()

		Expect(len(snapshot.Times)).To(Equal(2))
	})

	It("Creates a snapshot after a time period", func() {
		var controller = CreateAggregatorController(metrics.DefaultRegistry)

		var snapshot AggregatorSnapShot

		controller.Start()
		time.Sleep(2 * time.Second)
		controller.Stop()
		snapshot = controller.Snapshot()

		Expect(len(snapshot.Times)).To(Equal(3))
	})
})
