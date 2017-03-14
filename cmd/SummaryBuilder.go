package cmd

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/guzzlerio/corcel/core"
	"github.com/nsf/termbox-go"
)

//ConsoleSummaryBuilder ...
type ConsoleSummaryBuilder struct {
	writer io.Writer
	width  int
	height int
}

//NewConsoleSummaryBuilder ...
func NewConsoleSummaryBuilder(writer io.Writer) *ConsoleSummaryBuilder {

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	termbox.Close()

	return &ConsoleSummaryBuilder{
		writer: writer,
		width:  w,
		height: h,
	}
}

//Write ...
func (this *ConsoleSummaryBuilder) Write(summary core.ExecutionSummary) {

	this.top()
	this.line("Running Time", summary.RunningTime.String())
	this.line("Throughput", fmt.Sprintf("%-.0f req/s", summary.Throughput))
	this.line("Total Requests", fmt.Sprintf("%-.0f", summary.TotalRequests))
	this.line("Number of Errors", fmt.Sprintf("%-.0f", summary.TotalErrors))
	this.line("Availability", fmt.Sprintf("%-.4f%%", summary.Availability))
	this.line("Bytes Sent", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.Sent.Total))))
	this.line("Bytes Received", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.Received.Total))))
	this.line("Mean Response Time", fmt.Sprintf("%.4f ms", summary.ResponseTime.Mean))
	this.line("Min Response Time", fmt.Sprintf("%.4f ms", summary.ResponseTime.Min))
	this.line("Max Response Time", fmt.Sprintf("%.4f ms", summary.ResponseTime.Max))
	this.tail()
}

func (this *ConsoleSummaryBuilder) top() {
	fmt.Fprintln(this.writer, "╔═════════════════════════════════════════════════╗")
	fmt.Fprintln(this.writer, "║                     Summary                     ║")
	fmt.Fprintln(this.writer, "╠═════════════════════════════════════════════════╣")
}

func (this *ConsoleSummaryBuilder) tail() {
	fmt.Fprintln(this.writer, "╚═════════════════════════════════════════════════╝")
}

func (this *ConsoleSummaryBuilder) line(label string, value string) {
	data := fmt.Sprintf("%23s: %-22s", label, value)
	fmt.Fprintf(this.writer, "║ %s ║\n", data)
}
