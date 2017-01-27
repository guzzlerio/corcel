package inproc

import (
	"strings"
	"time"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/statistics"
	"github.com/rcrowley/go-metrics"
)

//Lock ...
//var Lock = &sync.Mutex{}

//Throughput ...
//var Throughput = 0

//ProcessEventsSubscribed ...
//var ProcessEventsSubscribed = 0

//ProcessEventsPublished ...
//var ProcessEventsPublished = 0

//NewGeneralExecutionResultProcessor ...
func NewGeneralExecutionResultProcessor() GeneralExecutionResultProcessor {
	return GeneralExecutionResultProcessor{}
}

//GeneralExecutionResultProcessor ...
type GeneralExecutionResultProcessor struct {
}

//Process ...
func (instance GeneralExecutionResultProcessor) Process(result core.ExecutionResult, registry metrics.Registry) {

	//TODO: Check whether this was duration urn before.  It looks like it was but not sure
	obj := result[core.DurationUrn.String()]
	timer := metrics.GetOrRegisterTimer(core.DurationUrn.Timer().String(), registry)
	timer.Update(obj.(time.Duration))

	timerH := metrics.GetOrRegisterHistogram(core.DurationUrn.Histogram().String(), registry, metrics.NewUniformSample(100))
	timerH.Update(int64(obj.(time.Duration)))

	throughput := metrics.GetOrRegisterMeter(core.ThroughputUrn.Meter().String(), registry)
	throughput.Mark(1)

	//	Throughput = Throughput + 1

	errors := metrics.GetOrRegisterMeter(core.ErrorUrn.Meter().String(), registry)

	if result[core.ErrorUrn.String()] != nil {
		var errorString string

		switch t := result[core.ErrorUrn.String()].(type) {
		case error:
			errorString = t.Error()
		case string:
			errorString = t
		}
		if !strings.Contains(errorString, "net/http: request canceled") {
			errors.Mark(1)
			statistics.IncrementCounter(registry, core.ErrorUrn.Counter().String(), 1)
		}
	}

	if result[core.BytesSentCountUrn.String()] != nil {
		bytesSentValue := int64(result[core.BytesSentCountUrn.String()].(int))

		urn := core.BytesSentCountUrn.Counter().String()
		statistics.IncrementCounter(registry, urn, bytesSentValue)

		bytesSent := metrics.GetOrRegisterHistogram(core.BytesSentCountUrn.Histogram().String(), registry, metrics.NewUniformSample(100))
		bytesSent.Update(bytesSentValue)
	}

	if result[core.BytesReceivedCountUrn.String()] != nil {
		bytesReceivedValue := int64(result[core.BytesReceivedCountUrn.String()].(int))

		statistics.IncrementCounter(registry, core.BytesReceivedCountUrn.Counter().String(), bytesReceivedValue)

		bytesReceived := metrics.GetOrRegisterHistogram(core.BytesReceivedCountUrn.Histogram().String(), registry, metrics.NewUniformSample(100))
		bytesReceived.Update(int64(result[core.BytesReceivedCountUrn.String()].(int)))
	}

	if result[core.AssertionsUrn.String()] != nil {
		assertionResults := result[core.AssertionsUrn.String()].([]core.AssertionResult)
		for _, result := range assertionResults {
			statistics.IncrementCounter(registry, core.AssertionsTotalUrn.Counter().String(), 1)

			if result[core.AssertionResultUrn.String()].(bool) == false {
				statistics.IncrementCounter(registry, core.AssertionsFailedUrn.Counter().String(), 1)
			}
		}
	}
}
