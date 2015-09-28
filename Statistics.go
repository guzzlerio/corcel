package main

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type Statistics struct {
	hBytesSent     metrics.Histogram
	hBytesReceived metrics.Histogram
	hResponseTime  metrics.Histogram
	mBytesSent     metrics.Meter
	mBytesReceived metrics.Meter
	start          time.Time
	mRequests      metrics.Meter
}

func CreateStatistics() *Statistics {
	sampleSize := 1024
	return &Statistics{
		hBytesSent:     metrics.NewHistogram(metrics.NewUniformSample(sampleSize)),
		hBytesReceived: metrics.NewHistogram(metrics.NewUniformSample(sampleSize)),
		hResponseTime:  metrics.NewHistogram(metrics.NewUniformSample(sampleSize)),
		mBytesSent:     metrics.NewMeter(),
		mBytesReceived: metrics.NewMeter(),
		mRequests: metrics.NewMeter(),
	}
}

func (instance *Statistics) Start() {
	instance.start = time.Now()
}

func (instance *Statistics) Request() {
	instance.mRequests.Mark(1)
}

func (instance *Statistics) BytesReceived(value int64) {
	instance.hBytesReceived.Update(value)
	instance.mBytesReceived.Mark(value)
}

func (instance *Statistics) BytesSent(value int64) {
	instance.hBytesSent.Update(value)
	instance.mBytesSent.Mark(value)
}

func (instance *Statistics) ResponseTime(value int64) {
	instance.hResponseTime.Update(value)
}

func (instance *Statistics) ExecutionOutput() ExecutionOutput {
	output := ExecutionOutput{
		Summary: ExecutionSummary{
			Requests : RequestsSummary{
				Rate : instance.mRequests.RateMean(),
			},
			RunningTime: float64(time.Since(instance.start) / time.Millisecond),
			ResponseTime: ResponseTimeStats{
				Sum:    instance.hResponseTime.Sum(),
				Max:    instance.hResponseTime.Max(),
				Mean:   instance.hResponseTime.Mean(),
				Min:    instance.hResponseTime.Min(),
				P50:    instance.hResponseTime.Percentile(50),
				P75:    instance.hResponseTime.Percentile(75),
				P95:    instance.hResponseTime.Percentile(95),
				P99:    instance.hResponseTime.Percentile(99),
				StdDev: instance.hResponseTime.StdDev(),
				Var:    instance.hResponseTime.Variance(),
			},
			Bytes: BytesSummary{
				Sent: BytesStats{
					Sum:    instance.hBytesSent.Sum(),
					Max:    instance.hBytesSent.Max(),
					Mean:   instance.hBytesSent.Mean(),
					Min:    instance.hBytesSent.Min(),
					P50:    instance.hBytesSent.Percentile(50),
					P75:    instance.hBytesSent.Percentile(75),
					P95:    instance.hBytesSent.Percentile(95),
					P99:    instance.hBytesSent.Percentile(99),
					StdDev: instance.hBytesSent.StdDev(),
					Var:    instance.hBytesSent.Variance(),
					Rate:   instance.mBytesSent.RateMean(),
				},
				Received: BytesStats{
					Sum:    instance.hBytesReceived.Sum(),
					Max:    instance.hBytesReceived.Max(),
					Mean:   instance.hBytesReceived.Mean(),
					Min:    instance.hBytesReceived.Min(),
					P50:    instance.hBytesReceived.Percentile(50),
					P75:    instance.hBytesReceived.Percentile(75),
					P95:    instance.hBytesReceived.Percentile(95),
					P99:    instance.hBytesReceived.Percentile(99),
					StdDev: instance.hBytesReceived.StdDev(),
					Var:    instance.hBytesReceived.Variance(),
					Rate:   instance.mBytesReceived.RateMean(),
				},
			},
		},
	}
	return output
}
