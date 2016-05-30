package main

import (
	"fmt"
	"net/http"

	"ci.guzzler.io/guzzler/corcel/statistics"
	"ci.guzzler.io/guzzler/corcel/test"
	"ci.guzzler.io/guzzler/corcel/utils"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run Summary", func() {

	FIt("Creates a summary file if one does not exist", func() {
		TestServer.Clear()
		TestServer.Use(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		}).For(rizo.RequestWithPath("/people"))

		planBuilder := test.NewYamlPlanBuilder()
		planBuilder.SetDuration("5s").
			CreateJob().
			CreateStep().
			ToExecuteAction(GetPathRequest(TestServer.CreateURL("/people")))

		err := ExecutePlanBuilder(planBuilder)
		fmt.Println(err)

		var executionOutput statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./output.yml", &executionOutput)

		var executionHistory statistics.AggregatorSnapShot
		utils.UnmarshalYamlFromFile("./history.yml", &executionHistory)

		for _, counters := range executionHistory.Counters {
			Expect(len(counters)).To(Equal(1))
		}
		for _, gauges := range executionHistory.Guages {
			Expect(len(gauges)).To(Equal(1))
		}
		for _, value := range executionHistory.Histograms {
			for _, subValue := range value {
				Expect(len(subValue)).To(Equal(1))
			}
		}
		for _, value := range executionHistory.Meters {
			for _, subValue := range value {
				Expect(len(subValue)).To(Equal(1))
			}
		}
		for _, value := range executionHistory.Timers {
			for _, subValue := range value {
				Expect(len(subValue)).To(Equal(1))
			}
		}
		Expect(len(executionHistory.Times)).To(Equal(1))

		Expect(err).To(BeNil())
		TestServer.Clear()
		/*
		 */
	})

})
