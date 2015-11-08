package main

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

//Statistics ...
type Statistics struct {
	hBytesSent     metrics.Histogram
	hBytesReceived metrics.Histogram
	hResponseTime  metrics.Histogram
	mBytesSent     metrics.Meter
	mBytesReceived metrics.Meter
	start          time.Time
	end            time.Time
	mRequests      metrics.Meter
	cErrors        metrics.Counter
	cTotal         metrics.Counter
}

//CreateStatistics ...
func CreateStatistics() *Statistics {
	sampleSize := 1024
	return &Statistics{
		hBytesSent:     metrics.NewHistogram(metrics.NewUniformSample(sampleSize)),
		hBytesReceived: metrics.NewHistogram(metrics.NewUniformSample(sampleSize)),
		hResponseTime:  metrics.NewHistogram(metrics.NewUniformSample(sampleSize)),
		mBytesSent:     metrics.NewMeter(),
		mBytesReceived: metrics.NewMeter(),
		mRequests:      metrics.NewMeter(),
		cErrors:        metrics.NewCounter(),
		cTotal:         metrics.NewCounter(),
	}
}

//Start ...
func (instance *Statistics) Start() {
	instance.start = time.Now()
}

//Stop ...
func (instance *Statistics) Stop() {
	instance.end = time.Now()
}

//Request ...
func (instance *Statistics) Request(err error) {
	instance.mRequests.Mark(1)
	if err != nil {
		instance.cErrors.Inc(1)
	}
	instance.cTotal.Inc(1)
}

//BytesReceived ...
func (instance *Statistics) BytesReceived(value int64) {
	instance.hBytesReceived.Update(value)
	instance.mBytesReceived.Mark(value)
}

//BytesSent ...
func (instance *Statistics) BytesSent(value int64) {
	instance.hBytesSent.Update(value)
	instance.mBytesSent.Mark(value)
}

//ResponseTime ...
func (instance *Statistics) ResponseTime(value int64) {
	instance.hResponseTime.Update(value)
}

//ExecutionOutput ...
func (instance *Statistics) ExecutionOutput() ExecutionOutput {
	//TODO: Separate out the construction of the following variables into methods
	//so that it is easier to manage.  Bit lengthy this ...
	output := ExecutionOutput{
		Summary: ExecutionSummary{
			Requests: RequestsSummary{
				Rate:         instance.mRequests.RateMean(),
				Errors:       instance.cErrors.Count(),
				Total:        instance.cTotal.Count(),
				Availability: 1 - (float64(instance.cErrors.Count()) / float64(instance.cTotal.Count())),
			},
			RunningTime: float64(instance.end.Sub(instance.start) / time.Millisecond),
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
