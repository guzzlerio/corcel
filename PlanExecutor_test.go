package main

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/cmd"
	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
	"ci.guzzler.io/guzzler/corcel/processor"
	. "ci.guzzler.io/guzzler/corcel/utils"

	"github.com/REAANDREW/rizo"
)

var _ = Describe("Plan Executor", func() {
	var list []string
	var file *os.File
	var server *rizo.RequestRecordingServer

	BeforeEach(func() {
		server = rizo.CreateRequestRecordingServer(5001)
		list = []string{
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/1")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/2")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/3")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/4")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/5")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/6")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/7")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/8")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/9")),
			fmt.Sprintf(`%s -X POST `, server.CreateURL("/10")),
		}
		file = CreateFileFromLines(list)
		server.Start()
	})

	AfterEach(func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
		server.Stop()
	})

	It("Does something", func() {
		configuration := config.DefaultConfig()
		configuration.FilePath = file.Name()
		bar := cmd.NewProgressBar(100, &configuration)
		stats := processor.CreateStatistics()

		executor := processor.PlanExecutor{
			Config: &configuration,
			Bar:    bar,
			Stats:  stats,
		}
		executor.Execute()

		Expect(len(server.Requests)).To(Equal(len(list)))
	})

	PIt("URL File", func() {})
	PIt("URL File with duration", func() {})
	PIt("URL File with random selection", func() {})
	PIt("URL File with more than one worker", func() {})
	PIt("URL File with wait time", func() {})
})
