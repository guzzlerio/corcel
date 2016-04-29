package http

import (
	"fmt"
	"time"

	"ci.guzzler.io/guzzler/corcel/core"

	"github.com/rcrowley/go-metrics"
)

//NewHTTPExecutionResultProcessor ...
func NewHTTPExecutionResultProcessor() HTTPExecutionResultProcessor {
	return HTTPExecutionResultProcessor{}
}

//HTTPExecutionResultProcessor ...
type HTTPExecutionResultProcessor struct {
}

//Process ...
func (instance HTTPExecutionResultProcessor) Process(result core.ExecutionResult, registry metrics.Registry) {
	for key, value := range result {
		switch key {
		case "http:request:error":
			meter := metrics.GetOrRegisterMeter("http:request:error", registry)
			meter.Mark(1)

			url := result["http:request:url"]

			byURLRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byURL:%s:", url))
			byURLmeter := metrics.GetOrRegisterMeter("http:request:error", byURLRegistry)
			byURLmeter.Mark(1)

		case "http:response:error":
			meter := metrics.GetOrRegisterMeter("http:response:error", registry)
			meter.Mark(1)

			url := result["http:request:url"]
			byURLRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byURL:%s:", url))
			byURLmeter := metrics.GetOrRegisterMeter("http:response:error", byURLRegistry)
			byURLmeter.Mark(1)

		case "http:request:bytes":
			url := result["http:request:url"]
			byURLRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byURL:%s:", url))
			byURLHistogram := metrics.GetOrRegisterHistogram("http:request:bytes", byURLRegistry, metrics.NewUniformSample(100))
			byURLHistogram.Update(int64(value.(int)))

		case "http:response:bytes":
			url := result["http:request:url"]
			byURLRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byURL:%s:", url))
			byURLHistogram := metrics.GetOrRegisterHistogram("http:response:bytes", byURLRegistry, metrics.NewUniformSample(100))
			byURLHistogram.Update(int64(value.(int)))

		case "http:response:status":
			statusCode := value.(int)
			url := result["http:request:url"]
			counter := metrics.GetOrRegisterCounter(fmt.Sprintf("http:response:status:%d", statusCode), registry)
			counter.Inc(1)

			byURLRegistry := metrics.NewPrefixedChildRegistry(registry, fmt.Sprintf("byURL:%s:", url))
			byURLCounter := metrics.GetOrRegisterCounter(fmt.Sprintf("http:response:status:%d", statusCode), byURLRegistry)
			byURLCounter.Inc(1)

			obj := result["action:duration"]
			timer := metrics.GetOrRegisterTimer(fmt.Sprintf("http:response:status:%d:duration", statusCode), registry)
			timer.Update(obj.(time.Duration))

			byURLTimer := metrics.GetOrRegisterTimer(fmt.Sprintf("http:response:status:%d:duration", statusCode), byURLRegistry)
			byURLTimer.Update(obj.(time.Duration))
		}
	}
}
