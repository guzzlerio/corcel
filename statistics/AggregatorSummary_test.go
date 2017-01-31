package statistics

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AggregatorSnapShot", func() {

	var (
		snapshot *AggregatorSnapShot
	)

	BeforeEach(func() {
		snapshot = NewAggregatorSnapShot()
	})

	It("handles an empty snapshot", func() {
		Ω(snapshot.CreateSummary()).ShouldNot(BeNil())
	})

	Describe("setting the RunningTime", func() {
		It("it calculates the correct duration", func() {
			start := int64(1483228800000)
			end := int64(1483229010000)
			duration := time.Unix(0, end).Sub(time.Unix(0, start))
			snapshot.updateTime(start)
			snapshot.updateTime(end)
			summary := snapshot.CreateSummary()
			Ω(summary.RunningTime).Should(Equal(duration))
		})

		It("does not permit a negative result", func() {
			//Invert the times from previous test
			end := int64(1483228800000)
			start := int64(1483229010000)
			snapshot.updateTime(start)
			snapshot.updateTime(end)
			summary := snapshot.CreateSummary()
			Ω(summary.RunningTime).Should(BeZero())
		})
	})
})
