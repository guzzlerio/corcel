package processor

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
	. "ci.guzzler.io/guzzler/corcel/utils"
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
	})

	AfterEach(func() {
		err := os.Remove(file.Name())
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
		//server.Stop()
	})

	It("URL File with duration", func() {
		start := time.Now()
		configuration.Duration = time.Duration(5 * time.Second)

		executor := CreatePlanExecutor(&configuration, bar)
		plan := CreatePlanFromURLList(&configuration)
		executor.Execute(plan)

		duration := time.Since(start)
		Expect(int(duration / time.Second)).To(Equal(5))
	})

	It("URL File with random selection", func() {
		configuration.Random = true
		executor := CreatePlanExecutor(&configuration, bar)

		tries := 50
		firstPaths := []string{}
		plan := CreatePlanFromURLList(&configuration)
		for i := 0; i < tries; i++ {
			executor.Execute(plan)
			if !ContainsString(firstPaths, TestServer.Requests[0].Request.URL.Path) {
				firstPaths = append(firstPaths, TestServer.Requests[0].Request.URL.Path)
			}
			TestServer.Clear()
		}
		Expect(len(firstPaths)).To(BeNumerically(">", 1))

	})

	It("URL File with more than one worker", func() {
		configuration.Workers = 5

		executor := CreatePlanExecutor(&configuration, bar)

		plan := CreatePlanFromURLList(&configuration)
		executor.Execute(plan)

		Expect(len(TestServer.Requests)).To(Equal(configuration.Workers * len(list)))
	})

	It("URL File with wait time", func() {
		waitTimeInMilliseconds := 200
		expectedTotalTimeInMilliseconds := len(list) * waitTimeInMilliseconds
		configuration.WaitTime = time.Duration(time.Duration(waitTimeInMilliseconds) * time.Millisecond)

		executor := CreatePlanExecutor(&configuration, bar)
		plan := CreatePlanFromURLList(&configuration)

		start := time.Now()
		executor.Execute(plan)

		duration := time.Since(start)
		Expect(int(duration / time.Second)).To(Equal(expectedTotalTimeInMilliseconds / 1000))
	})
})
