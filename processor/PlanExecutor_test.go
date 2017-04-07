package processor

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	metrics "github.com/rcrowley/go-metrics"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/statistics"
	. "github.com/guzzlerio/corcel/utils"
)

//NullProgressBar ...
type NullProgressBar struct {
}

func (instance NullProgressBar) Set(progress int) error {
	return nil
}

func TestPlan_Executor(t *testing.T) {
	BeforeTest()
	defer AfterTest()
	Convey("Plan Executor", t, func() {
		var list []string
		var file *os.File
		var configuration *config.Configuration
		var bar ProgressBar
		var aggregator statistics.AggregatorInterfaceToRenameLater
		var channel chan core.ExecutionResult

		func() {
			list = []string{
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/1")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/2")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/3")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/4")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/5")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/6")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/7")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/8")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/9")),
				fmt.Sprintf(`%s -X POST `, TestServer.CreateURL("/10")),
			}
			configuration = config.DefaultConfig()
			file = CreateFileFromLines(list)
			configuration.FilePath = file.Name()
			bar = NullProgressBar{}
			aggregator = statistics.NewAggregator(metrics.DefaultRegistry)
			channel = make(chan core.ExecutionResult)

			go func() {
				for range channel {
				}
			}()
		}()

		defer func() {
			err := os.Remove(file.Name())
			if err != nil {
				logger.Log.Printf("Error removing file %v", err)
			}
			close(channel)
		}()

		Convey("URL File with duration", func() {
			start := time.Now()
			configuration.Duration = time.Duration(5 * time.Second)

			executor := CreatePlanExecutor(configuration, bar, core.CreateRegistry(), aggregator, channel)
			executor.Execute()

			duration := time.Since(start)
			So(int(duration/time.Second), ShouldEqual, 5)
		})

		Convey("URL File with random selection", func() {
			configuration.Random = true

			tries := 50
			firstPaths := []string{}
			for i := 0; i < tries; i++ {
				executor := CreatePlanExecutor(configuration, bar, core.CreateRegistry(), aggregator, channel)
				executor.Execute()
				if !ContainsString(firstPaths, TestServer.Requests[0].Request.URL.Path) {
					firstPaths = append(firstPaths, TestServer.Requests[0].Request.URL.Path)
				}
				TestServer.Clear()
			}
			So(len(firstPaths), ShouldBeGreaterThan, 1)

		})

		Convey("URL File with more than one worker", func() {
			configuration.Workers = 5

			executor := CreatePlanExecutor(configuration, bar, core.CreateRegistry(), aggregator, channel)

			executor.Execute()

			So(len(TestServer.Requests), ShouldEqual, configuration.Workers*len(list))
		})

		SkipConvey("URL File with wait time", func() {
			waitTimeInMilliseconds := 200
			expectedTotalTimeInMilliseconds := len(list) * waitTimeInMilliseconds
			configuration.WaitTime = time.Duration(time.Duration(waitTimeInMilliseconds) * time.Millisecond)

			executor := CreatePlanExecutor(configuration, bar, core.CreateRegistry(), aggregator, channel)

			start := time.Now()
			executor.Execute()

			duration := time.Since(start)
			So(int(duration/time.Second), ShouldEqual, expectedTotalTimeInMilliseconds/1000)
		})
	})
}
