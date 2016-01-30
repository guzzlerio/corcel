package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/cmd"
	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/processor"
	. "ci.guzzler.io/guzzler/corcel/utils"
)

var _ = Describe("Plan Executor", func() {
	var list []string
	var file *os.File
	var stats *processor.Statistics
	//var server *rizo.RequestRecordingServer
	var configuration config.Configuration
	var bar processor.ProgressBar

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
		bar = cmd.NewProgressBar(100, &configuration)
		stats = processor.CreateStatistics()
		//server.Start()
	})

	AfterEach(func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
		//server.Stop()
	})

	It("URL File", func() {
		executor := processor.PlanExecutor{
			Config: &configuration,
			Bar:    bar,
			Stats:  stats,
		}
		executor.Execute()

		Expect(len(TestServer.Requests)).To(Equal(len(list)))
	})

	It("URL File updates the Statistics", func() {
		executor := processor.PlanExecutor{
			Config: &configuration,
			Bar:    bar,
			Stats:  stats,
		}
		executor.Execute()
		output := stats.ExecutionOutput()
		Expect(output.Summary.Requests.Total).To(Equal(int64(len(list))))
	})

	It("URL File with duration", func() {
		start := time.Now()
		configuration.Duration = time.Duration(5 * time.Second)

		executor := processor.PlanExecutor{
			Config: &configuration,
			Bar:    bar,
			Stats:  stats,
		}

		executor.Execute()
		stats.ExecutionOutput()

		duration := time.Since(start)
		Expect(int(duration / time.Second)).To(Equal(5))
	})
	PIt("URL File with random selection", func() {})
	PIt("URL File with more than one worker", func() {})
	PIt("URL File with wait time", func() {})
})
