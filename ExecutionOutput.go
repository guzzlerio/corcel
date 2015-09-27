package main

type ResponseTimeStats struct{
	Sum    int64 `yaml:"sum"`
	Max    int64 `yaml:"max"`
	Mean   float64 `yaml:"mean"`
	Min    int64 `yaml:"min"`
	P50    float64 `yaml:"p50"`
	P75    float64 `yaml:"p75"`
	P95    float64 `yaml:"p95"`
	P99    float64 `yaml:"p99"`
	StdDev float64 `yaml:"stddev"`
	Var    float64 `yaml:"var"`
}

type BytesStats struct {
	Sum    int64 `yaml:"sum"`
	Max    int64 `yaml:"max"`
	Mean   float64 `yaml:"mean"`
	Min    int64 `yaml:"min"`
	P50    float64 `yaml:"p50"`
	P75    float64 `yaml:"p75"`
	P95    float64 `yaml:"p95"`
	P99    float64 `yaml:"p99"`
	StdDev float64 `yaml:"stddev"`
	Var    float64 `yaml:"var"`
	Rate   float64 `yaml:"rate"`
}

type BytesSummary struct {
	Sent     BytesStats `yaml:"sent"`
	Received BytesStats `yaml:"received"`
}

type ExecutionSummary struct {
	Bytes BytesSummary `yaml:"bytes"`
	ResponseTime ResponseTimeStats `yaml:"responseTime"`
	RunningTime float64 `yaml:"runningTime"`
}

type ExecutionOutput struct {
	Summary ExecutionSummary `yaml:"summary"`
}
