package statistics

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("AggregatorSnapShot", func() {

	var (
		targetSnapShot  *AggregatorSnapShot
		subjectSnapShot *AggregatorSnapShot
	)

	BeforeEach(func() {
		targetSnapShot = NewAggregatorSnapShot()
		subjectSnapShot = NewAggregatorSnapShot()
	})

	It("Updates an empty snap shot counters", func() {

		key := "counter:1"
		value := int64(101)

		timeStamp := time.Now().UnixNano()
		subjectSnapShot.updateCounter(key, value)
		subjectSnapShot.updateTime(timeStamp)

		targetSnapShot.Update(*subjectSnapShot)

		Expect(len(targetSnapShot.Times)).To(Equal(1))
		Expect(targetSnapShot.Times[0]).To(Equal(subjectSnapShot.Times[0]))
		Expect(targetSnapShot.Counters[key]).ToNot(BeNil())
		Expect(targetSnapShot.Counters[key]).To(Equal(subjectSnapShot.Counters[key]))
	})

	It("Updates an empty snap shot gauges", func() {
		key := "gauge:1"
		value := float64(102.1)

		timeStamp := time.Now().UnixNano()
		subjectSnapShot.updateGauge(key, value)
		subjectSnapShot.updateTime(timeStamp)

		targetSnapShot.Update(*subjectSnapShot)

		Expect(len(targetSnapShot.Times)).To(Equal(1))
		Expect(targetSnapShot.Times[0]).To(Equal(subjectSnapShot.Times[0]))
		Expect(targetSnapShot.Gauges[key]).ToNot(BeNil())
		Expect(targetSnapShot.Gauges[key]).To(Equal(subjectSnapShot.Gauges[key]))
	})

	It("Updates an empty snap shot histograms", func() {
		key := "h:key"
		subKey := "h:subkey"
		value := float64(103.1)

		timeStamp := time.Now().UnixNano()
		subjectSnapShot.updateHistogram(key, subKey, value)
		subjectSnapShot.updateTime(timeStamp)

		targetSnapShot.Update(*subjectSnapShot)

		Expect(len(targetSnapShot.Times)).To(Equal(1))
		Expect(targetSnapShot.Times[0]).To(Equal(subjectSnapShot.Times[0]))
		Expect(targetSnapShot.Histograms[key]).ToNot(BeNil())
		Expect(targetSnapShot.Histograms[key][subKey]).ToNot(BeNil())
		Expect(targetSnapShot.Histograms[key][subKey]).To(Equal(subjectSnapShot.Histograms[key][subKey]))
	})
})
