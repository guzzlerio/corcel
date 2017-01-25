package core

import (
	"fmt"
	"io"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/ghodss/yaml"
	"github.com/nsf/termbox-go"
)

//SummaryBuilder ...
type SummaryBuilder interface {
	Write(summary ExecutionSummary)
}

func NewSummaryBuilder(format string) SummaryBuilder {
	switch format {
	case "json":
		return &JSONSummaryBuilder{}
	default:
		return NewConsoleSummaryBuilder()
	}
}

type ConsoleSummaryBuilder struct {
	writer io.Writer
	width  int
	height int
}

func NewConsoleSummaryBuilder() *ConsoleSummaryBuilder {

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	termbox.Close()

	return &ConsoleSummaryBuilder{
		writer: os.Stdout,
		width:  w,
		height: h,
	}
}

func (i *ConsoleSummaryBuilder) Write(summary ExecutionSummary) {

	fmt.Fprintf(i.writer, "w: %v h: %v\n", i.width, i.height)
	i.top()
	i.line("Running Time", summary.RunningTime)
	i.line("Throughput", fmt.Sprintf("%-.0f req/s", summary.Throughput))
	i.line("Total Requests", fmt.Sprintf("%-.0f", summary.TotalRequests))
	i.line("Number of Errors", fmt.Sprintf("%-.0f", summary.TotalErrors))
	i.line("Availability", fmt.Sprintf("%-.4f%%", summary.Availability))
	i.line("Bytes Sent", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.Sent.Total))))
	i.line("Bytes Received", fmt.Sprintf("%v", humanize.Bytes(uint64(summary.Bytes.Received.Total))))
	i.line("Mean Response Time", fmt.Sprintf("%.4f ms", summary.MeanResponseTime))
	i.line("Min Response Time", fmt.Sprintf("%.4f ms", summary.MinResponseTime))
	i.line("Max Response Time", fmt.Sprintf("%.4f ms", summary.MaxResponseTime))
	i.tail()
}

func (i *ConsoleSummaryBuilder) top() {
	fmt.Fprintln(i.writer, "╔═════════════════════════════════════════════════╗")
	fmt.Fprintln(i.writer, "║                     Summary                     ║")
	fmt.Fprintln(i.writer, "╠═════════════════════════════════════════════════╣")
}

func (i *ConsoleSummaryBuilder) tail() {
	fmt.Fprintln(i.writer, "╚═════════════════════════════════════════════════╝")
}

func (i *ConsoleSummaryBuilder) line(label string, value string) {
	data := fmt.Sprintf("%23s: %-22s", label, value)
	// fmt.Fprintf(i.writer, "data len: %v\n", len(label+": "+value))
	fmt.Fprintf(i.writer, "║ %s ║\n", data)
}

type JSONSummaryBuilder struct {
}

func (this *JSONSummaryBuilder) Write(summary ExecutionSummary) {
	y, _ := yaml.Marshal(summary)
	jsonData, _ := yaml.YAMLToJSON(y)
	fmt.Println(string(jsonData))
}
