package statistics

import (
	"time"

	"github.com/guzzlerio/corcel/core"
)

//CreateSummary ...
func CreateSummary(snapshot AggregatorSnapShot) core.ExecutionSummary {

	lastTime := time.Unix(0, snapshot.Times[len(snapshot.Times)-1])
	firstTime := time.Unix(0, snapshot.Times[0])
	duration := lastTime.Sub(firstTime)

	counts := snapshot.Meters[core.ThroughputUrn.Meter().String()]["count"]
	count := counts[len(counts)-1]

	errors := snapshot.Meters[core.ErrorUrn.Meter().String()]["count"]
	errorCount := errors[len(errors)-1]

	rates := snapshot.Meters[core.ThroughputUrn.Meter().String()]["rateMean"]
	rate := rates[len(rates)-1]

	var availability float64
	if errorCount > 0 {
		availability = (1 - (float64(errorCount) / float64(count))) * 100
	} else {
		availability = 100
	}

	var totalAssertionsCount = int64(0)
	var totalAssertionFailuresCount = int64(0)
	var totalReceived = int64(0)
	var totalSent = int64(0)

	bytes := core.ByteSummary{}

	bytesSent := snapshot.Counters[core.BytesSentCountUrn.Counter().String()]

	if bytesSent != nil {
		for _, value := range bytesSent {
			totalSent += value
		}
	}

	bytesReceived := snapshot.Counters[core.BytesReceivedCountUrn.Counter().String()]
	if bytesReceived != nil {
		for _, value := range bytesReceived {
			totalReceived += value
		}
	}

	bytesReceivedHistogram := snapshot.Histograms[core.BytesReceivedCountUrn.Histogram().String()]
	if bytesReceivedHistogram != nil {
		bytes.Received = core.ByteStat{
			Min:   int64FromHistogram(bytesReceivedHistogram["min"]),
			Max:   int64FromHistogram(bytesReceivedHistogram["max"]),
			Mean:  int64FromHistogram(bytesReceivedHistogram["mean"]),
			Total: totalReceived,
		}
	}

	bytesSentHistogram := snapshot.Histograms[core.BytesSentCountUrn.Histogram().String()]
	if bytesSentHistogram != nil {
		bytes.Sent = core.ByteStat{
			Min:   int64FromHistogram(bytesSentHistogram["min"]),
			Max:   int64FromHistogram(bytesSentHistogram["max"]),
			Mean:  int64FromHistogram(bytesSentHistogram["mean"]),
			Total: totalSent,
		}
	}

	totalAssertions := snapshot.Counters[core.AssertionsTotalUrn.Counter().String()]
	if totalAssertions != nil {
		totalAssertionsCount = totalAssertions[len(totalAssertions)-1]
	}

	totalAssertionsFailed := snapshot.Counters[core.AssertionsFailedUrn.Counter().String()]
	if totalAssertionsFailed != nil {
		totalAssertionFailuresCount = totalAssertionsFailed[len(totalAssertionsFailed)-1]
	}
	responseMeanTimes := snapshot.Timers[core.DurationUrn.Timer().String()]["mean"]
	responseMeanTime := responseMeanTimes[len(responseMeanTimes)-1]

	responseMinTimes := snapshot.Timers[core.DurationUrn.Timer().String()]["min"]
	responseMinTime := responseMinTimes[len(responseMinTimes)-1]

	responseMaxTimes := snapshot.Timers[core.DurationUrn.Timer().String()]["max"]
	responseMaxTime := responseMaxTimes[len(responseMaxTimes)-1]

	return core.ExecutionSummary{
		RunningTime:   duration,
		TotalRequests: count,
		TotalErrors:   errorCount,
		Availability:  availability,
		Throughput:    rate,
		ResponseTime: core.ResponseTimeStat{
			Mean: responseMeanTime,
			Min:  responseMinTime,
			Max:  responseMaxTime,
		},
		TotalAssertions:        totalAssertionsCount,
		TotalAssertionFailures: totalAssertionFailuresCount,
		Bytes: bytes,
	}
}

func int64FromHistogram(b []int64) int64 {
	return b[len(b)-1]
}
