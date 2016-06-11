package statistics

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("AggregatorSnapShot", func() {

	It("Updates an empty snap shot counters", func() {

		key := "counter:1"
		value := int64(101)
		timeStamp := time.Now().UnixNano()
		targetSnapShot := NewAggregatorSnapShot()

		subjectSnapShot := NewAggregatorSnapShot()
		subjectSnapShot.updateCounter(key, value)
		subjectSnapShot.updateTime(timeStamp)

		targetSnapShot.Update(*subjectSnapShot)

		Expect(len(targetSnapShot.Times)).To(Equal(1))
		Expect(targetSnapShot.Times[0]).To(Equal(subjectSnapShot.Times[0]))

	})
})
