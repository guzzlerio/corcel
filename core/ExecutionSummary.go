package core

type MinMaxMeanTotalInt struct {
	Min   int64 `json:"min"`
	Max   int64 `json:"max"`
	Mean  int64 `json:"mean"`
	Total int64 `json:"total"`
}

//ByteSummary ...
type ByteSummary struct {
	Received MinMaxMeanTotalInt `json:"received"`
	Sent     MinMaxMeanTotalInt `json:"sent"`
}

//ExecutionSummary ...
type ExecutionSummary struct {
	TotalRequests          float64     `json:"totalRequests"`
	TotalErrors            float64     `json:"totalErrors"`
	Availability           float64     `json:"availability"`
	RunningTime            string      `json:"runningTime"`
	Throughput             float64     `json:"throughput"`
	MeanResponseTime       float64     `json:"meanResponseTime"`
	MinResponseTime        float64     `json:"minResponseTime"`
	MaxResponseTime        float64     `json:"maxResponseTime"`
	TotalAssertions        int64       `json:"totalAssertions"`
	TotalAssertionFailures int64       `json:"totalAssertionFailures"`
	Bytes                  ByteSummary `json:"bytes"`
}
