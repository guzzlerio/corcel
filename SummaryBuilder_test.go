package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/guzzlerio/corcel/cmd"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/serialisation/json"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSummary_Builders(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("Summary Builders", t, func() {
		var (
			builder core.SummaryBuilder
			summary core.ExecutionSummary
			writer  *bytes.Buffer
		)
		func() {
			writer = new(bytes.Buffer)
			duration, _ := time.ParseDuration("1.003242751s")
			summary = core.ExecutionSummary{
				Availability: 100,
				Bytes: core.ByteSummary{
					Received: core.ByteStat{
						Max:   120,
						Mean:  120,
						Min:   120,
						Total: 303480,
					},
					Sent: core.ByteStat{
						Max:   308,
						Mean:  136,
						Min:   59,
						Total: 343864,
					},
				},
				ResponseTime: core.ResponseTimeStat{
					Max:  31,
					Mean: 6.7081714,
					Min:  0,
				},
				RunningTime:            duration,
				Throughput:             2544.8113,
				TotalAssertionFailures: 0,
				TotalAssertions:        0,
				TotalErrors:            0,
				TotalRequests:          2540,
			}
		}()

		Convey("ConsoleSummaryBuilder", func() {
			func() {
				builder = cmd.NewConsoleSummaryBuilder(writer)
			}()

			Convey("writes the expected table to the console", func() {
				builder.Write(summary)
				So(writer.String(), ShouldContainSubstring, `╔═════════════════════════════════════════════════╗
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
╚═════════════════════════════════════════════════╝`)
			})
		})

		Convey("JSONSummaryBuilder", func() {
			func() {
				builder = &json.JSONSummaryBuilder{writer}
			}()

			Convey("writes the expected table to the console", func() {
				builder.Write(summary)
				So(writer.String(), json.ShouldMatchJson, `
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
  "runningTime": 1003242751,
  "throughput": 2544.8113,
  "totalAssertionFailures": 0,
  "totalAssertions": 0,
  "totalErrors": 0,
  "totalRequests": 2540
}
`)
			})
		})

		Convey("YAMLSummaryBuilder", func() {
			func() {
				builder = &yaml.YAMLSummaryBuilder{writer}
			}()

			Convey("writes the expected table to the console", func() {
				builder.Write(summary)
				So(writer.String(), yaml.ShouldMatchYaml, `
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
runningTime: 1003242751
throughput: 2544.8113
totalAssertionFailures: 0
totalAssertions: 0
totalErrors: 0
totalRequests: 2540
`)
			})
		})
	})
}
