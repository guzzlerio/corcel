package statistics

import (
	"time"

	"github.com/guzzlerio/corcel/core"
)

//ByteSummary ...
type ByteSummary struct {
	MinReceived   int64
	MaxReceived   int64
	MeanReceived  int64
	MinSent       int64
	MaxSent       int64
	MeanSent      int64
	TotalSent     int64
	TotalReceived int64
}

//ExecutionSummary ...
type ExecutionSummary struct {
	TotalRequests          float64
	TotalErrors            float64
	Availability           float64
	RunningTime            string
	Throughput             float64
	MeanResponseTime       float64
	MinResponseTime        float64
	MaxResponseTime        float64
	TotalAssertions        int64
	TotalAssertionFailures int64
	Bytes                  ByteSummary
}

//CreateSummary ...
func CreateSummary(snapshot AggregatorSnapShot) ExecutionSummary {

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

	bytes := ByteSummary{}

	bytesSent := snapshot.Counters[core.BytesSentCountUrn.Counter().String()]

	if bytesSent != nil {
		//bytesSentCount = bytesSent[len(bytesSent)-1]

		for _, value := range bytesSent {
			bytes.TotalSent += value
		}
	}

	bytesReceived := snapshot.Counters[core.BytesReceivedCountUrn.Counter().String()]
	if bytesReceived != nil {
		//bytesReceivedCount = bytesReceived[len(bytesReceived)-1]
		for _, value := range bytesReceived {
			bytes.TotalReceived += value
		}
	}

	bytesReceivedHistogram := snapshot.Histograms[core.BytesReceivedCountUrn.Histogram().String()]
	if bytesReceivedHistogram != nil {
		maxBytes := bytesReceivedHistogram["max"]
		bytes.MaxReceived = maxBytes[len(maxBytes)-1]

		minBytes := bytesReceivedHistogram["min"]
		bytes.MinReceived = minBytes[len(minBytes)-1]
	}

	bytesSentHistogram := snapshot.Histograms[core.BytesSentCountUrn.Histogram().String()]
	if bytesSentHistogram != nil {
		maxBytes := bytesSentHistogram["max"]
		bytes.MaxSent = maxBytes[len(maxBytes)-1]

		minBytes := bytesSentHistogram["min"]
		bytes.MinSent = minBytes[len(minBytes)-1]
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

	return ExecutionSummary{
		RunningTime:            duration.String(),
		TotalRequests:          count,
		TotalErrors:            errorCount,
		Availability:           availability,
		Throughput:             rate,
		MeanResponseTime:       responseMeanTime,
		MinResponseTime:        responseMinTime,
		MaxResponseTime:        responseMaxTime,
		TotalAssertions:        totalAssertionsCount,
		TotalAssertionFailures: totalAssertionFailuresCount,
		Bytes: bytes,
	}
}
