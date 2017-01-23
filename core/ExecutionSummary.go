package core

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
