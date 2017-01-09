package http

import (
	"time"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/statistics"

	"github.com/rcrowley/go-metrics"
)

//NewExecutionResultProcessor ...
func NewExecutionResultProcessor() ExecutionResultProcessor {
	return ExecutionResultProcessor{}
}

//ExecutionResultProcessor ...
type ExecutionResultProcessor struct {
}

//Process ...
func (instance ExecutionResultProcessor) Process(result core.ExecutionResult, registry metrics.Registry) {

	for key, value := range result {
		switch key {
		case RequestErrorUrn.String():
			meter := metrics.GetOrRegisterMeter(RequestErrorUrn.Meter().String(), registry)
			meter.Mark(1)

		case ResponseErrorUrn.String():
			meter := metrics.GetOrRegisterMeter(ResponseErrorUrn.Meter().String(), registry)
			meter.Mark(1)

		case core.BytesSentCountUrn.String():
			byURLHistogram := metrics.GetOrRegisterHistogram(RequestBytesUrn.Histogram().String(), registry, metrics.NewUniformSample(100))
			byURLHistogram.Update(int64(value.(int)))

		case core.BytesReceivedCountUrn.String():
			byURLHistogram := metrics.GetOrRegisterHistogram(ResponseBytesUrn.Histogram().String(), registry, metrics.NewUniformSample(100))
			byURLHistogram.Update(int64(value.(int)))

		case ResponseStatusUrn.String():
			statusCode := value.(int)
			statistics.IncrementCounter(registry, ResponseStatusUrn.Name(statusCode).Counter().String(), 1)

			obj := result[core.DurationUrn.String()]
			timer := metrics.GetOrRegisterTimer(ResponseStatusUrn.Name(statusCode).Timer().String(), registry)
			timer.Update(obj.(time.Duration))
		}
	}
}
