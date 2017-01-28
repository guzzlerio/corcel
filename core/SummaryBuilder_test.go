package core

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Summary Builders", func() {
	var (
		builder SummaryBuilder
		summary ExecutionSummary
		writer  *bytes.Buffer
	)
	BeforeEach(func() {
		writer = new(bytes.Buffer)
		summary = ExecutionSummary{
			Availability: 100,
			Bytes: ByteSummary{
				Received: ByteStat{
					Max:   120,
					Mean:  120,
					Min:   120,
					Total: 303480,
				},
				Sent: ByteStat{
					Max:   308,
					Mean:  136,
					Min:   59,
					Total: 343864,
				},
			},
			ResponseTime: ResponseTimeStat{
				Max:  31,
				Mean: 6.7081714,
				Min:  0,
			},
			RunningTime:            "1.003242751s",
			Throughput:             2544.8113,
			TotalAssertionFailures: 0,
			TotalAssertions:        0,
			TotalErrors:            0,
			TotalRequests:          2540,
		}
	})

	Describe("ConsoleSummaryBuilder", func() {
		BeforeEach(func() {
			builder = NewConsoleSummaryBuilder(writer)
		})

		It("writes the expected table to the console", func() {
			builder.Write(summary)
			Ω(writer.String()).Should(MatchRegexp(`╔═════════════════════════════════════════════════╗
║                     Summary                     ║
╠═════════════════════════════════════════════════╣
║            Running Time: 1.003242751s           ║
║              Throughput: 2545 req/s             ║
║          Total Requests: 2540                   ║
║        Number of Errors: 0                      ║
║            Availability: 100.0000%              ║
║              Bytes Sent: 344 kB                 ║
║          Bytes Received: 304 kB                 ║
║      Mean Response Time: 6.7082 ms              ║
║       Min Response Time: 0.0000 ms              ║
║       Max Response Time: 31.0000 ms             ║
╚═════════════════════════════════════════════════╝`))
		})
	})

	Describe("JSONSummaryBuilder", func() {
		BeforeEach(func() {
			builder = &JSONSummaryBuilder{writer}
		})

		It("writes the expected table to the console", func() {
			builder.Write(summary)
			Ω(writer.String()).Should(MatchJSON(`
{
"availability": 100,
  "bytes": {
    "received": {
      "max": 120,
      "mean": 120,
      "min": 120,
      "total": 303480
    },
    "sent": {
      "max": 308,
      "mean": 136,
      "min": 59,
      "total": 343864
    }
  },
  "responseTime": {
    "max": 31,
    "mean": 6.7081714,
    "min": 0
  },
  "runningTime": "1.003242751s",
  "throughput": 2544.8113,
  "totalAssertionFailures": 0,
  "totalAssertions": 0,
  "totalErrors": 0,
  "totalRequests": 2540
}
`))
		})
	})

	Describe("YAMLSummaryBuilder", func() {
		BeforeEach(func() {
			builder = &YAMLSummaryBuilder{writer}
		})

		It("writes the expected table to the console", func() {
			builder.Write(summary)
			Ω(writer.String()).Should(MatchYAML(`
availability: 100
bytes:
  received:
    max: 120
    mean: 120
    min: 120
    total: 303480
  sent:
    max: 308
    mean: 136
    min: 59
    total: 343864
responseTime:
  max: 31
  mean: 6.7081714
  min: 0
runningTime: 1.003242751s
throughput: 2544.8113
totalAssertionFailures: 0
totalAssertions: 0
totalErrors: 0
totalRequests: 2540
`))
		})
	})
})
