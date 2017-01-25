package statistics

import (
	"regexp"
	"strconv"
	"time"
)

//ExecutionSummaryCliParser ...
type ExecutionSummaryCliParser struct {
}

//Parse ...
func (instance ExecutionSummaryCliParser) Parse(data string) ExecutionSummary {
	return ExecutionSummary{
		RunningTime:   parseRunningTime(data),
		Throughput:    parseThroughput(data),
		TotalRequests: parseTotalRequests(data),
		TotalErrors:   parseTotalErrors(data),
		Availability:  parseAvailability(data),
		Bytes: ByteSummary{
			TotalSent:     parseTotalBytesSent(data),
			TotalReceived: parseTotalBytesReceived(data),
		},
		MeanResponseTime: parseMeanResponseTime(data),
		MinResponseTime:  parseMinResponseTime(data),
		MaxResponseTime:  parseMaxResponseTime(data),
	}
}

func checkMatchString(match []string) string {
	if len(match) != 2 {
		return ""
	}
	return match[1]
}

func checkMatchFloat(match []string) float64 {
	if len(match) != 2 {
		return float64(-1)
	}
	value, _ := strconv.ParseFloat(match[1], 64)
	return value
}

func checkMatchInt64(match []string) int64 {
	if len(match) != 2 {
		return int64(-1)
	}
	value, _ := strconv.Atoi(match[1])
	return int64(value)
}

func checkMatchTime(match []string) time.Duration {
	if len(match) != 2 {
		return time.Duration(-1)
	}
	value, _ := time.ParseDuration(match[1])
	return value
}

func parseRunningTime(data string) time.Duration {
	var runningTime = regexp.MustCompile(`Running Time: ([\d\w\.]+)`)
	return checkMatchTime(runningTime.FindStringSubmatch(data))
}

func parseThroughput(data string) float64 {
	var thoughput = regexp.MustCompile(`Throughput: ([\d]+)`)
	return checkMatchFloat(thoughput.FindStringSubmatch(data))
}

func parseTotalRequests(data string) float64 {
	var totalRequests = regexp.MustCompile(`Total Requests: ([\d]+)`)
	return checkMatchFloat(totalRequests.FindStringSubmatch(data))
}

func parseTotalErrors(data string) float64 {
	var totalErrors = regexp.MustCompile(`Number of Errors: ([\d]+)`)
	return checkMatchFloat(totalErrors.FindStringSubmatch(data))
}

func parseAvailability(data string) float64 {
	var availability = regexp.MustCompile(`Availability: ([\d\.\d]+)%`)
	return checkMatchFloat(availability.FindStringSubmatch(data))
}

func parseTotalBytesSent(data string) int64 {
	var totalBytesSent = regexp.MustCompile(`Bytes Sent: ([\d]+)`)
	return checkMatchInt64(totalBytesSent.FindStringSubmatch(data))
}

func parseTotalBytesReceived(data string) int64 {
	var totalBytesReceived = regexp.MustCompile(`Bytes Received: ([\d]+)`)
	return checkMatchInt64(totalBytesReceived.FindStringSubmatch(data))
}

func parseMeanResponseTime(data string) float64 {
	var meanResponseTime = regexp.MustCompile(`Mean Response Time: ([\d\.\d]+)`)
	return checkMatchFloat(meanResponseTime.FindStringSubmatch(data))
}

func parseMinResponseTime(data string) float64 {
	var minResponseTime = regexp.MustCompile(`Min Response Time: ([\d\.\d]+)`)
	return checkMatchFloat(minResponseTime.FindStringSubmatch(data))
}

func parseMaxResponseTime(data string) float64 {
	var maxResponseTime = regexp.MustCompile(`Max Response Time: ([\d\.\d]+)`)
	return checkMatchFloat(maxResponseTime.FindStringSubmatch(data))
}

//CreateExecutionSummaryCliParser ...
func CreateExecutionSummaryCliParser() ExecutionSummaryCliParser {
	return ExecutionSummaryCliParser{}
}
