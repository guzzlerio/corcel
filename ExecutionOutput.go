package main

import (
	"fmt"
)

//ResponseTimeStats ...
type ResponseTimeStats struct {
	Sum    int64   `yaml:"sum"`
	Max    int64   `yaml:"max"`
	Mean   float64 `yaml:"mean"`
	Min    int64   `yaml:"min"`
	P50    float64 `yaml:"p50"`
	P75    float64 `yaml:"p75"`
	P95    float64 `yaml:"p95"`
	P99    float64 `yaml:"p99"`
	StdDev float64 `yaml:"stddev"`
	Var    float64 `yaml:"var"`
}

//BytesStats ...
type BytesStats struct {
	Sum    int64   `yaml:"sum"`
	Max    int64   `yaml:"max"`
	Mean   float64 `yaml:"mean"`
	Min    int64   `yaml:"min"`
	P50    float64 `yaml:"p50"`
	P75    float64 `yaml:"p75"`
	P95    float64 `yaml:"p95"`
	P99    float64 `yaml:"p99"`
	StdDev float64 `yaml:"stddev"`
	Var    float64 `yaml:"var"`
	Rate   float64 `yaml:"rate"`
}

//BytesSummary ...
type BytesSummary struct {
	Sent     BytesStats `yaml:"sent"`
	Received BytesStats `yaml:"received"`
}

//RequestsSummary ...
type RequestsSummary struct {
	Rate         float64 `yaml:"rate"`
	Errors       int64   `yaml:"errors"`
	Total        int64   `yaml:"total"`
	Availability float64 `yaml:"availability"`
}

//ExecutionSummary ...
type ExecutionSummary struct {
	Bytes        BytesSummary      `yaml:"bytes"`
	ResponseTime ResponseTimeStats `yaml:"responseTime"`
	RunningTime  float64           `yaml:"runningTime"`
	Requests     RequestsSummary   `yaml:"requests"`
}

//ExecutionOutput ...
type ExecutionOutput struct {
	Summary ExecutionSummary `yaml:"summary"`
}

type ExecutionOutputConsoleWriter struct {
	Output ExecutionOutput
}

func (w *ExecutionOutputConsoleWriter) Write() {
	top()
	line("Running Time", fmt.Sprintf("%g s", w.Output.Summary.RunningTime/1000))
	line("Throughput", fmt.Sprintf("%-v req/s", int64(w.Output.Summary.Requests.Rate)))
	line("Total Requests", fmt.Sprintf("%v", w.Output.Summary.Requests.Total))
	line("Number of Errors", fmt.Sprintf("%v", w.Output.Summary.Requests.Errors))
	line("Availability", fmt.Sprintf("%.4v%%", w.Output.Summary.Requests.Availability*100))
	line("Bytes Sent", fmt.Sprintf("%v", w.Output.Summary.Bytes.Sent.Sum))
	line("Bytes Received", fmt.Sprintf("%v", w.Output.Summary.Bytes.Received.Sum))
	if w.Output.Summary.ResponseTime.Mean > 0 {
		line("Mean Response Time", fmt.Sprintf("%.4v ms", w.Output.Summary.ResponseTime.Mean))
	} else {
		line("Mean Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Mean))
	}

	line("Min Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Min))
	line("Max Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Max))
	tail()
}

func top() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                           Summary                                 ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════════╣")
}

func tail() {
	fmt.Println("╚═══════════════════════════════════════════════════════════════════╝")
}

func line(label string, value string) {
	fmt.Printf("║ %20s: %-43s ║\n", label, value)
}
