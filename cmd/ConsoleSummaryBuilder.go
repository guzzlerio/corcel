package cmd

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/statistics"
	"io"
)

type ConsoleSummaryBuilder struct {
	writer io.Writer
}

func NewConsoleSummaryBuilder() *ConsoleSummaryBuilder {
	return ConsoleSummaryBuilder{
		writer: os.Stdout,
	}
}

func (i *core.SummaryBuilder) Write(summary core.ExecutionSummary) {

	top()
	line("Running Time", summary.RunningTime)
	line("Throughput", fmt.Sprintf("%-.0f req/s", summary.Throughput))
	line("Total Requests", fmt.Sprintf("%-.0f", summary.TotalRequests))
	line("Number of Errors", fmt.Sprintf("%-.0f", summary.TotalErrors))
	line("Availability", fmt.Sprintf("%-.4f%%", summary.Availability))
	line("Bytes Sent", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.TotalSent))))
	line("Bytes Received", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.TotalReceived))))
	line("Mean Response Time", fmt.Sprintf("%.4f ms", summary.MeanResponseTime))
	line("Min Response Time", fmt.Sprintf("%.4f ms", summary.MinResponseTime))
	line("Max Response Time", fmt.Sprintf("%.4f ms", summary.MaxResponseTime))
	tail()
}

func (i *core.SummaryBuilder) top(writer io.Writer) {
	fmt.Fprintln(i.writer, "╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(i.writer, "║                           Summary                                 ║")
	fmt.Fprintln(i.writer, "╠═══════════════════════════════════════════════════════════════════╣")
}

func (i *core.SummaryBuilder) tail(writer io.Writer) {
	fmt.Fprintln(i.writer, "╚═══════════════════════════════════════════════════════════════════╝")
}

func (i *core.SummaryBuilder) line(writer io.Writer, label string, value string) {
	fmt.Fprintf(i.writer, "║ %20s: %-43s ║\n", label, value)
}
