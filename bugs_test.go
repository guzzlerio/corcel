package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/guzzlerio/corcel/errormanager"
	"github.com/guzzlerio/corcel/logger"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	. "github.com/guzzlerio/corcel/utils"
)

var _ = Describe("Bugs replication", func() {

	BeforeEach(func() {
		err := os.Remove("./output.yml")
		if err != nil {
			logger.Log.Printf("Error removing file %v", err)
		}
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	It("Error when running a simple run with POST and PUT #18", func() {
		numberOfWorkers := 2
		list := []string{
			fmt.Sprintf(`%s -X POST -d '{"name": "bob"}' -H "Content-type: application/json"`, URLForTestServer("/success")),
			fmt.Sprintf(`%s -X PUT -d '{"id": 1,"name": "bob junior"}' -H "Content-type: application/json"`, URLForTestServer("/success")),
			fmt.Sprintf(`%s?id=1 -X GET -H "Content-type: application/json"`, URLForTestServer("/success")),
			fmt.Sprintf(`%s?id=1 -X DELETE -H "Content-type: application/json"`, URLForTestServer("/success")),
		}

		SutExecute(list[:1], "--random", "--summary", "--workers", strconv.Itoa(numberOfWorkers))

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(summary.TotalRequests).To(Equal(float64(2)))
	})

	PIt("Error when too many workers specified causing too many open files #23", func() {
		numberOfWorkers := 100000000000
		list := []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/error")),
		}

		output, err := InvokeCorcel(list, "--workers", strconv.Itoa(numberOfWorkers))

		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring("Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers"))
	})

	It("Error non-http url in the urls file causes a run time exception #21", func() {
		list := []string{
			fmt.Sprintf(`-Something`),
		}

		output, err := InvokeCorcel(list)
		Expect(err).ToNot(BeNil())
		Expect(string(output)).To(ContainSubstring(errormanager.LogMessageVaidURLs))
	})

	It("Issue #49 - Corcel not cancelling on-going requests once the test is due to finish", func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		})
		TestServer.Use(factory)
		list := []string{
			fmt.Sprintf(`%s -X POST `, URLForTestServer("/something")),
		}

		SutExecute(list, "--duration", "1s")

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		runningTime, _ := time.ParseDuration(summary.RunningTime)
		Expect(math.Floor(runningTime.Seconds())).To(Equal(float64(1)))
	})

	It("Issue - Should write out panics to a log file and not std out", func() {
		planBuilder := yaml.NewPlanBuilder()

		planBuilder.
			SetIterations(1).
			CreateJob().
			CreateStep().
			ToExecuteAction(planBuilder.IPanicAction().Build())

		output, err := ExecutePlanBuilder(planBuilder)
		Expect(err).ToNot(BeNil())

		Expect(string(output)).To(ContainSubstring("An unexpected error has occurred.  The error has been logged to /tmp/"))

		//Ensure that the file which was generated contains the error which caused the panic
		r, _ := regexp.Compile(`/tmp/[\w\d-]+`)
		var location = r.FindString(string(output))
		Expect(location).ToNot(Equal(""))
		data, err := ioutil.ReadFile(location)
		Expect(err).To(BeNil())
		Expect(string(data)).To(ContainSubstring("IPanicAction has caused this panic"))
	})
})
