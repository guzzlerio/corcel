package statistics

import (
	"math"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionSummaryCliParser", func() {

	Describe("Empty Data Path", func() {
		var data = ""
		var parser = CreateExecutionSummaryCliParser()
		var executionSummary ExecutionSummary

		BeforeEach(func() {
			executionSummary = parser.Parse(data)
		})

		It("Parses Running Time", func() {
			Expect(executionSummary.RunningTime).To(Equal(time.Duration(-1)))
		})

		It("Parses Throughput", func() {
			Expect(executionSummary.Throughput).To(Equal(float64(-1)))
		})

		It("Parses Total Requests", func() {
			Expect(executionSummary.TotalRequests).To(Equal(float64(-1)))
		})

		It("Parses Number of Errors", func() {
			Expect(executionSummary.TotalErrors).To(Equal(float64(-1)))
		})

		It("Parses Availability", func() {
			Expect(executionSummary.Availability).To(Equal(float64(-1)))
		})

		It("Parses Bytes Sent", func() {
			Expect(executionSummary.Bytes.TotalSent).To(Equal(int64(-1)))
		})

		It("Parses Bytes Received", func() {
			Expect(executionSummary.Bytes.TotalReceived).To(Equal(int64(-1)))
		})

		It("Parses Mean Response Time", func() {
			Expect(executionSummary.MeanResponseTime).To(Equal(float64(-1)))
		})

		It("Parses Min Response Time", func() {
			Expect(executionSummary.MinResponseTime).To(Equal(float64(-1)))
		})

		It("Parses Max Response Time", func() {
			Expect(executionSummary.MaxResponseTime).To(Equal(float64(-1)))
		})
	})

	Describe("Happy Path", func() {
		var data = `
╔═══════════════════════════════════════════════════════════════════╗
║                           Summary                                 ║
╠═══════════════════════════════════════════════════════════════════╣
║         Running Time: 3.093643s                                   ║
║           Throughput: 232288 req/s                                ║
║       Total Requests: 1                                           ║
║     Number of Errors: 0                                           ║
║         Availability: 100.0000%                                   ║
║           Bytes Sent: 74 B                                        ║
║       Bytes Received: 144 B                                       ║
║   Mean Response Time: 2.0000 ms                                   ║
║    Min Response Time: 2.0000 ms                                   ║
║    Max Response Time: 2.0000 ms                                   ║
╚═══════════════════════════════════════════════════════════════════╝
	`
		var parser = CreateExecutionSummaryCliParser()
		var executionSummary ExecutionSummary

		BeforeEach(func() {
			executionSummary = parser.Parse(data)
		})

		It("Parses Running Time", func() {
			Expect(math.Floor(executionSummary.RunningTime.Seconds())).To(Equal(float64(3)))
		})

		It("Parses Throughput", func() {
			Expect(executionSummary.Throughput).To(Equal(float64(232288)))
		})

		It("Parses Total Requests", func() {
			Expect(executionSummary.TotalRequests).To(Equal(float64(1)))
		})

		It("Parses Number of Errors", func() {
			Expect(executionSummary.TotalErrors).To(Equal(float64(0)))
		})

		It("Parses Availability", func() {
			Expect(executionSummary.Availability).To(Equal(float64(100)))
		})

		It("Parses Bytes Sent", func() {
			Expect(executionSummary.Bytes.TotalSent).To(Equal(int64(74)))
		})

		It("Parses Bytes Received", func() {
			Expect(executionSummary.Bytes.TotalReceived).To(Equal(int64(144)))
		})

		It("Parses Mean Response Time", func() {
			Expect(executionSummary.MeanResponseTime).To(Equal(float64(2)))
		})

		It("Parses Min Response Time", func() {
			Expect(executionSummary.MinResponseTime).To(Equal(float64(2)))
		})

		It("Parses Max Response Time", func() {
			Expect(executionSummary.MaxResponseTime).To(Equal(float64(2)))
		})
	})
})
