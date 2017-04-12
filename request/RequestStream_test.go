package request

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	. "github.com/guzzlerio/corcel/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestStream(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("RequestStream", t, func() {
		var (
			list     []string
			reader   *Reader
			iterator RequestStream
		)

		func() {
			list = []string{
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/1")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/2")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/3")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/4")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/5")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/6")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/7")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/8")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/9")),
				fmt.Sprintf(`%s -X POST `, URLForTestServer("/10")),
			}
			file := CreateFileFromLines(list)
			check(file.Close())
			reader = NewRequestReader(file.Name())
		}()

		Convey("Sequential Request Stream", func() {
			func() {
				iterator = NewSequentialRequestStream(reader)
			}()

			Convey("executes in a sequential order", func() {
				actual := []*http.Request{}
				for iterator.HasNext() {
					req, _ := iterator.Next()
					actual = append(actual, req)
				}
				So(len(actual), ShouldEqual, len(list))
			})

			Convey("calculates Size", func() {
				So(iterator.Size(), ShouldEqual, len(list))
			})

			Convey("calculates Progress", func() {
				for i := 0; i < (len(list) / 2); i++ {
					_, _ = iterator.Next()
				}
				So(iterator.Progress(), ShouldEqual, 50)
			})
		})

		Convey("Timebased RequestStream", func() {
			var duration time.Duration
			var lock = &sync.Mutex{}

			func() {
				lock.Lock()
				iterator = NewSequentialRequestStream(reader)
				duration = time.Duration(1 * time.Second)
				iterator = NewTimeBasedRequestStream(iterator, duration)
				lock.Unlock()
			}()

			Convey("executes for the duration", func() {
				actual := Time(func() {
					for iterator.HasNext() {
						_, err := iterator.Next()
						check(err)
					}
				})

				var min = duration - (1 * time.Millisecond)
				var max = duration + (50 * time.Millisecond)

				So(DurationIsBetween(actual, min, max), ShouldEqual, true)
			})

			Convey("calculates Progress", func() {
				go func() {
					lock.Lock()
					for iterator.HasNext() {
						_, _ = iterator.Next()
					}
					lock.Unlock()
				}()
				time.Sleep((100 * time.Millisecond))
				So(iterator.Progress(), ShouldBeGreaterThanOrEqualTo, 10)
				So(iterator.Progress(), ShouldBeLessThan, 20)
			})

			Convey("calculates Size", func() {
				So(iterator.Size(), ShouldEqual, int(duration.Nanoseconds()))
			})
		})

		Convey("Random RequestStream", func() {
			Convey("calculates Size", func() {
				stream := NewRandomRequestStream(reader)
				So(stream.Size(), ShouldEqual, len(list))
			})

			Convey("executes in a random order", func() {
				requestSet1 := []*http.Request{}
				requestSet2 := []*http.Request{}

				stream1 := NewRandomRequestStream(reader)
				for stream1.HasNext() {
					req, _ := stream1.Next()
					requestSet1 = append(requestSet1, req)
				}

				stream2 := NewRandomRequestStream(reader)
				for stream2.HasNext() {
					req, _ := stream2.Next()
					requestSet2 = append(requestSet2, req)
				}

				So(len(requestSet1), ShouldEqual, len(list))
				So(len(requestSet2), ShouldEqual, len(list))
				So(ConcatRequestPaths(requestSet1), ShouldNotResemble, ConcatRequestPaths(requestSet2))
			})

			Convey("from Reader with size of 1", func() {
				list = []string{
					fmt.Sprintf(`%s -X POST `, URLForTestServer("/1")),
				}
				file := CreateFileFromLines(list)
				check(file.Close())
				reader = NewRequestReader(file.Name())

				requestSet := []*http.Request{}

				stream := NewRandomRequestStream(reader)
				for stream.HasNext() {
					req, _ := stream.Next()
					requestSet = append(requestSet, req)
				}
				So(len(requestSet), ShouldEqual, len(list))
			})
		})

	})
}
