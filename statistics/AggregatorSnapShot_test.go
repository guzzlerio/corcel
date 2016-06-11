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

	It("Updates an empty snap shot meters", func() {
		key := "m:key"
		subKey := "m:subkey"
		value := float64(104.1)

		timeStamp := time.Now().UnixNano()
		subjectSnapShot.updateMeter(key, subKey, value)
		subjectSnapShot.updateTime(timeStamp)

		targetSnapShot.Update(*subjectSnapShot)

		Expect(len(targetSnapShot.Times)).To(Equal(1))
		Expect(targetSnapShot.Times[0]).To(Equal(subjectSnapShot.Times[0]))
		Expect(targetSnapShot.Meters[key]).ToNot(BeNil())
		Expect(targetSnapShot.Meters[key][subKey]).ToNot(BeNil())
		Expect(targetSnapShot.Meters[key][subKey]).To(Equal(subjectSnapShot.Meters[key][subKey]))
	})

	It("Updates an empty snap shot timers", func() {
		key := "t:key"
		subKey := "t:subkey"
		value := float64(105.1)

		timeStamp := time.Now().UnixNano()
		subjectSnapShot.updateTimer(key, subKey, value)
		subjectSnapShot.updateTime(timeStamp)

		targetSnapShot.Update(*subjectSnapShot)

		Expect(len(targetSnapShot.Times)).To(Equal(1))
		Expect(targetSnapShot.Times[0]).To(Equal(subjectSnapShot.Times[0]))
		Expect(targetSnapShot.Timers[key]).ToNot(BeNil())
		Expect(targetSnapShot.Timers[key][subKey]).ToNot(BeNil())
		Expect(targetSnapShot.Timers[key][subKey]).To(Equal(subjectSnapShot.Timers[key][subKey]))
	})

	Context("Only updates with the last values of the subject snap shot for", func() {

		It("Counters", func() {
			count := 10
			key := "counter:1"
			targetSnapShot.updateCounter(key, int64(0))
			targetSnapShot.updateTime(time.Now().UnixNano())
			for i := 0; i < count; i++ {
				subjectSnapShot.updateCounter(key, int64(i+1))
				subjectSnapShot.updateTime(time.Now().UnixNano())
			}
			targetSnapShot.Update(*subjectSnapShot)

			Expect(len(targetSnapShot.Times)).To(Equal(2))
			Expect(targetSnapShot.Counters[key][1]).To(Equal(subjectSnapShot.Counters[key][count-1]))
		})

		It("Gauges", func() {
			count := 10
			key := "gauge:1"
			targetSnapShot.updateGauge(key, float64(0))
			targetSnapShot.updateTime(time.Now().UnixNano())
			for i := 0; i < count; i++ {
				subjectSnapShot.updateGauge(key, float64(i+1))
				subjectSnapShot.updateTime(time.Now().UnixNano())
			}
			targetSnapShot.Update(*subjectSnapShot)

			Expect(len(targetSnapShot.Times)).To(Equal(2))
			Expect(targetSnapShot.Gauges[key][1]).To(Equal(subjectSnapShot.Gauges[key][count-1]))
		})
	})

})
