package core

import "time"

//ByteStat ...
type ByteStat struct {
	Min   int64 `json:"min"`
	Max   int64 `json:"max"`
	Mean  int64 `json:"mean"`
	Total int64 `json:"total"`
}

//ByteSummary ...
type ByteSummary struct {
	Received ByteStat `json:"received"`
	Sent     ByteStat `json:"sent"`
}

//ResponseTimeStat ...
type ResponseTimeStat struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Mean float64 `json:"mean"`
}

//ExecutionSummary ...
type ExecutionSummary struct {
	TotalRequests          float64          `json:"totalRequests"`
	TotalErrors            float64          `json:"totalErrors"`
	Availability           float64          `json:"availability"`
	RunningTime            time.Duration    `json:"runningTime"`
	Throughput             float64          `json:"throughput"`
	ResponseTime           ResponseTimeStat `json:"responseTime"`
	TotalAssertions        int64            `json:"totalAssertions"`
	TotalAssertionFailures int64            `json:"totalAssertionFailures"`
	Bytes                  ByteSummary      `json:"bytes"`
}
