package statistics

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregatorSnapShot(t *testing.T) {
	Convey("AggregatorSnapShot", t, func() {

		var (
			targetSnapShot  *AggregatorSnapShot
			subjectSnapShot *AggregatorSnapShot
		)

		func() {
			targetSnapShot = NewAggregatorSnapShot()
			subjectSnapShot = NewAggregatorSnapShot()
		}()

		Convey("Updates an empty snap shot counters", func() {

			key := "counter:1"
			value := int64(101)

			timeStamp := time.Now().UnixNano()
			subjectSnapShot.updateCounter(key, value)
			subjectSnapShot.updateTime(timeStamp)

			targetSnapShot.Update(*subjectSnapShot)

			So(len(targetSnapShot.Times), ShouldEqual, 1)
			So(targetSnapShot.Times[0], ShouldEqual, subjectSnapShot.Times[0])
			So(targetSnapShot.Counters[key], ShouldNotBeNil)
			So(targetSnapShot.Counters[key], ShouldResemble, subjectSnapShot.Counters[key])
		})

		Convey("Updates an empty snap shot gauges", func() {
			key := "gauge:1"
			value := float64(102.1)

			timeStamp := time.Now().UnixNano()
			subjectSnapShot.updateGauge(key, value)
			subjectSnapShot.updateTime(timeStamp)

			targetSnapShot.Update(*subjectSnapShot)

			So(len(targetSnapShot.Times), ShouldEqual, 1)
			So(targetSnapShot.Times[0], ShouldEqual, subjectSnapShot.Times[0])
			So(targetSnapShot.Gauges[key], ShouldNotBeNil)
			So(targetSnapShot.Gauges[key], ShouldResemble, subjectSnapShot.Gauges[key])
		})

		Convey("Updates an empty snap shot histograms", func() {
			key := "h:key"
			subKey := "h:subkey"
			value := int64(103)

			timeStamp := time.Now().UnixNano()
			subjectSnapShot.updateHistogram(key, subKey, value)
			subjectSnapShot.updateTime(timeStamp)

			targetSnapShot.Update(*subjectSnapShot)

			So(len(targetSnapShot.Times), ShouldEqual, 1)
			So(targetSnapShot.Times[0], ShouldEqual, subjectSnapShot.Times[0])
			So(targetSnapShot.Histograms[key], ShouldNotBeNil)
			So(targetSnapShot.Histograms[key][subKey], ShouldNotBeNil)
			So(targetSnapShot.Histograms[key][subKey], ShouldResemble, subjectSnapShot.Histograms[key][subKey])
		})

		Convey("Updates an empty snap shot meters", func() {
			key := "m:key"
			subKey := "m:subkey"
			value := float64(104.1)

			timeStamp := time.Now().UnixNano()
			subjectSnapShot.updateMeter(key, subKey, value)
			subjectSnapShot.updateTime(timeStamp)

			targetSnapShot.Update(*subjectSnapShot)

			So(len(targetSnapShot.Times), ShouldEqual, 1)
			So(targetSnapShot.Times[0], ShouldEqual, subjectSnapShot.Times[0])
			So(targetSnapShot.Meters[key], ShouldNotBeNil)
			So(targetSnapShot.Meters[key][subKey], ShouldNotBeNil)
			So(targetSnapShot.Meters[key][subKey], ShouldResemble, subjectSnapShot.Meters[key][subKey])
		})

		Convey("Updates an empty snap shot timers", func() {
			key := "t:key"
			subKey := "t:subkey"
			value := float64(105.1)

			timeStamp := time.Now().UnixNano()
			subjectSnapShot.updateTimer(key, subKey, value)
			subjectSnapShot.updateTime(timeStamp)

			targetSnapShot.Update(*subjectSnapShot)

			So(len(targetSnapShot.Times), ShouldEqual, 1)
			So(targetSnapShot.Times[0], ShouldEqual, subjectSnapShot.Times[0])
			So(targetSnapShot.Timers[key], ShouldNotBeNil)
			So(targetSnapShot.Timers[key][subKey], ShouldNotBeNil)
			So(targetSnapShot.Timers[key][subKey], ShouldResemble, subjectSnapShot.Timers[key][subKey])
		})

		Convey("Only updates with the last values of the subject snap shot for", func() {

			Convey("Counters", func() {
				count := 10
				key := "counter:1"
				targetSnapShot.updateCounter(key, int64(0))
				targetSnapShot.updateTime(time.Now().UnixNano())
				for i := 0; i < count; i++ {
					subjectSnapShot.updateCounter(key, int64(i+1))
					subjectSnapShot.updateTime(time.Now().UnixNano())
				}
				targetSnapShot.Update(*subjectSnapShot)

				So(len(targetSnapShot.Times), ShouldEqual, 2)
				So(targetSnapShot.Counters[key][1], ShouldEqual, subjectSnapShot.Counters[key][count-1])
			})

			Convey("Gauges", func() {
				count := 10
				key := "gauge:1"
				targetSnapShot.updateGauge(key, float64(0))
				targetSnapShot.updateTime(time.Now().UnixNano())
				for i := 0; i < count; i++ {
					subjectSnapShot.updateGauge(key, float64(i+1))
					subjectSnapShot.updateTime(time.Now().UnixNano())
				}
				targetSnapShot.Update(*subjectSnapShot)

				So(len(targetSnapShot.Times), ShouldEqual, 2)
				So(targetSnapShot.Gauges[key][1], ShouldEqual, subjectSnapShot.Gauges[key][count-1])
			})

			Convey("Histograms", func() {
				count := 10
				key := "h:key"
				subKey := "h:subkey"
				targetSnapShot.updateHistogram(key, subKey, int64(0))
				targetSnapShot.updateTime(time.Now().UnixNano())
				for i := 0; i < count; i++ {
					subjectSnapShot.updateHistogram(key, subKey, int64(i+1))
					subjectSnapShot.updateTime(time.Now().UnixNano())
				}
				targetSnapShot.Update(*subjectSnapShot)

				So(len(targetSnapShot.Times), ShouldEqual, 2)
				So(targetSnapShot.Histograms[key][subKey][1], ShouldEqual, subjectSnapShot.Histograms[key][subKey][count-1])
			})

			Convey("Meters", func() {
				count := 10
				key := "m:key"
				subKey := "m:subkey"
				targetSnapShot.updateMeter(key, subKey, float64(0))
				targetSnapShot.updateTime(time.Now().UnixNano())
				for i := 0; i < count; i++ {
					subjectSnapShot.updateMeter(key, subKey, float64(i+1))
					subjectSnapShot.updateTime(time.Now().UnixNano())
				}
				targetSnapShot.Update(*subjectSnapShot)

				So(len(targetSnapShot.Times), ShouldEqual, 2)
				So(targetSnapShot.Meters[key][subKey][1], ShouldEqual, subjectSnapShot.Meters[key][subKey][count-1])
			})

			Convey("Timers", func() {
				count := 10
				key := "t:key"
				subKey := "t:subkey"
				targetSnapShot.updateTimer(key, subKey, float64(0))
				targetSnapShot.updateTime(time.Now().UnixNano())
				for i := 0; i < count; i++ {
					subjectSnapShot.updateTimer(key, subKey, float64(i+1))
					subjectSnapShot.updateTime(time.Now().UnixNano())
				}
				targetSnapShot.Update(*subjectSnapShot)

				So(len(targetSnapShot.Times), ShouldEqual, 2)
				So(targetSnapShot.Timers[key][subKey][1], ShouldEqual, subjectSnapShot.Timers[key][subKey][count-1])
			})
		})

	})
}
