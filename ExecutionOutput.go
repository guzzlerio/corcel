package main

import (
	"fmt"
	"io"
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

//ExecutionOutputWriter ...
type ExecutionOutputWriter struct {
	Output ExecutionOutput
}

//Write ...
func (w *ExecutionOutputWriter) Write(writer io.Writer) {
	top(writer)
	line(writer, "Running Time", fmt.Sprintf("%g s", w.Output.Summary.RunningTime/1000))
	line(writer, "Throughput", fmt.Sprintf("%-v req/s", int64(w.Output.Summary.Requests.Rate)))
	line(writer, "Total Requests", fmt.Sprintf("%v", w.Output.Summary.Requests.Total))
	line(writer, "Number of Errors", fmt.Sprintf("%v", w.Output.Summary.Requests.Errors))
	line(writer, "Availability", fmt.Sprintf("%.4v%%", w.Output.Summary.Requests.Availability*100))
	line(writer, "Bytes Sent", fmt.Sprintf("%v", w.Output.Summary.Bytes.Sent.Sum))
	line(writer, "Bytes Received", fmt.Sprintf("%v", w.Output.Summary.Bytes.Received.Sum))
	if w.Output.Summary.ResponseTime.Mean > 0 {
		line(writer, "Mean Response Time", fmt.Sprintf("%.4v ms", w.Output.Summary.ResponseTime.Mean))
	} else {
		line(writer, "Mean Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Mean))
	}

	line(writer, "Min Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Min))
	line(writer, "Max Response Time", fmt.Sprintf("%v ms", w.Output.Summary.ResponseTime.Max))
	tail(writer)
}

func top(writer io.Writer) {
	fmt.Fprintln(writer, "╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(writer, "║                           Summary                                 ║")
	fmt.Fprintln(writer, "╠═══════════════════════════════════════════════════════════════════╣")
}

func tail(writer io.Writer) {
	fmt.Fprintln(writer, "╚═══════════════════════════════════════════════════════════════════╝")
}

func line(writer io.Writer, label string, value string) {
	fmt.Fprintf(writer, "║ %20s: %-43s ║\n", label, value)
}
