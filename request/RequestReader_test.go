package request

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/guzzlerio/corcel/global"
	. "github.com/guzzlerio/corcel/utils"
)

func TestRequestReader(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("RequestReader", t, func() {

		var list []string
		var reader *Reader

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
			err := file.Close()
			if err != nil {
				panic(err)
			}
			reader = NewRequestReader(file.Name())
		}()

		Convey("Single reader iterates over lines in a file", func() {
			requests := []*http.Request{}
			stream := NewSequentialRequestStream(reader)
			for stream.HasNext() {
				req, err := stream.Next()
				check(err)
				requests = append(requests, req)
			}
			So(len(requests), ShouldEqual, len(list))
		})

		for _, numberOfWorkers := range global.NumberOfWorkersToTest {
			Convey(fmt.Sprintf("Multiple readers totalling %v iterates over lines in a file", numberOfWorkers), func() {
				var wg sync.WaitGroup
				var mutex = &sync.Mutex{}
				wg.Add(numberOfWorkers)
				requests := []*http.Request{}
				for i := 0; i < numberOfWorkers; i++ {
					go func() {
						stream := NewSequentialRequestStream(reader)
						for stream.HasNext() {
							mutex.Lock()
							req, err := stream.Next()
							check(err)
							requests = append(requests, req)
							mutex.Unlock()
						}
						wg.Done()
					}()
				}
				wg.Wait()
				So(len(requests), ShouldEqual, len(list)*numberOfWorkers)
			})
		}
	})
}
