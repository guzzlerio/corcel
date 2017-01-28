package processor

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("Plan Executor", func() {
	var list []string
	var file *os.File
	var configuration config.Configuration
	var bar ProgressBar
	var aggregator statistics.AggregatorInterfaceToRenameLater
	var channel chan core.ExecutionResult

	BeforeEach(func() {
		//server = rizo.CreateRequestRecordingServer(5001)
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
	})

	AfterEach(func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
		//server.Stop()
		close(channel)
	})

	It("URL File with duration", func() {
		start := time.Now()
		configuration.Duration = time.Duration(5 * time.Second)

		executor := CreatePlanExecutor(&configuration, bar, core.CreateRegistry(), aggregator, channel)
		executor.Execute()

		duration := time.Since(start)
		Expect(int(duration / time.Second)).To(Equal(5))
	})

	It("URL File with random selection", func() {
		configuration.Random = true

		tries := 50
		firstPaths := []string{}
		for i := 0; i < tries; i++ {
			executor := CreatePlanExecutor(&configuration, bar, core.CreateRegistry(), aggregator, channel)
			executor.Execute()
			if !ContainsString(firstPaths, TestServer.Requests[0].Request.URL.Path) {
				firstPaths = append(firstPaths, TestServer.Requests[0].Request.URL.Path)
			}
			TestServer.Clear()
		}
		Expect(len(firstPaths)).To(BeNumerically(">", 1))

	})

	It("URL File with more than one worker", func() {
		configuration.Workers = 5

		executor := CreatePlanExecutor(&configuration, bar, core.CreateRegistry(), aggregator, channel)

		executor.Execute()

		Expect(len(TestServer.Requests)).To(Equal(configuration.Workers * len(list)))
	})

	PIt("URL File with wait time", func() {
		waitTimeInMilliseconds := 200
		expectedTotalTimeInMilliseconds := len(list) * waitTimeInMilliseconds
		configuration.WaitTime = time.Duration(time.Duration(waitTimeInMilliseconds) * time.Millisecond)

		executor := CreatePlanExecutor(&configuration, bar, core.CreateRegistry(), aggregator, channel)

		start := time.Now()
		executor.Execute()

		duration := time.Since(start)
		Expect(int(duration / time.Second)).To(Equal(expectedTotalTimeInMilliseconds / 1000))
	})
})
