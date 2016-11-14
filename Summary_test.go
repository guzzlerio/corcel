package main

import (
	"net/http"
	"os"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/guzzlerio/corcel/utils"
	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run Summary", func() {

	var runJob = func() {
		TestServer.Use(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		}).For(rizo.RequestWithPath("/people"))

		planBuilder := yaml.NewPlanBuilder()
		planBuilder.SetDuration("1s").
			CreateJob().
			CreateStep().
			ToExecuteAction(GetHTTPRequestAction(TestServer.CreateURL("/people")))

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())
	}

	var assertSizeOfSummary = func(snapshot statistics.AggregatorSnapShot, size int) {
		for _, counters := range snapshot.Counters {
			Expect(len(counters)).To(Equal(size))
		}
		for _, gauges := range snapshot.Gauges {
			Expect(len(gauges)).To(Equal(size))
		}
		for _, value := range snapshot.Histograms {
			for _, subValue := range value {
				Expect(len(subValue)).To(Equal(size))
			}
		}
		for _, value := range snapshot.Meters {
			for _, subValue := range value {
				Expect(len(subValue)).To(Equal(size))
			}
		}
		for _, value := range snapshot.Timers {
			for _, subValue := range value {
				Expect(len(subValue)).To(Equal(size))
			}
		}
		Expect(len(snapshot.Times)).To(Equal(size))
	}

	BeforeEach(func() {
		_ = os.Remove("./history.yml")
		TestServer.Clear()
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	It("Creates a summary file if one does not exist", func() {
		runJob()
		var executionHistory statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./history.yml", &executionHistory)

		assertSizeOfSummary(executionHistory, 1)
	})

	It("Adds to a summary file if one already exists", func() {
		runJob()
		runJob()

		var executionHistory statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./history.yml", &executionHistory)

		assertSizeOfSummary(executionHistory, 2)
	})

})
