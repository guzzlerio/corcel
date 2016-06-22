package inproc

import (
	"strings"
	"time"

	"ci.guzzler.io/guzzler/corcel/core"
	"github.com/rcrowley/go-metrics"
)

//NewGeneralExecutionResultProcessor ...
func NewGeneralExecutionResultProcessor() GeneralExecutionResultProcessor {
	return GeneralExecutionResultProcessor{}
}

//GeneralExecutionResultProcessor ...
type GeneralExecutionResultProcessor struct {
}

//Process ...
func (instance GeneralExecutionResultProcessor) Process(result core.ExecutionResult, registry metrics.Registry) {
	obj := result["action:duration"]
	timer := metrics.GetOrRegisterTimer("action:duration", registry)
	timer.Update(obj.(time.Duration))

	throughput := metrics.GetOrRegisterMeter("action:throughput", registry)
	throughput.Mark(1)

	errors := metrics.GetOrRegisterMeter("action:error", registry)
	if result["action:error"] != nil {
		var errorString string

		switch t := result["action:error"].(type) {
		case error:
			errorString = t.Error()
		case string:
			errorString = t
		}
		if !strings.Contains(errorString, "net/http: request canceled") {
			errors.Mark(1)
		}
	}

	if result["action:bytes:sent"] != nil {
		bytesSentValue := int64(result["action:bytes:sent"].(int))

		bytesSentCounter := metrics.GetOrRegisterCounter("counter:action:bytes:sent", registry)
		bytesSentCounter.Inc(bytesSentValue)

		bytesSent := metrics.GetOrRegisterHistogram("histogram:action:bytes:sent", registry, metrics.NewUniformSample(100))
		bytesSent.Update(bytesSentValue)
	}

	if result["action:bytes:received"] != nil {
		bytesReceivedValue := int64(result["action:bytes:received"].(int))

		bytesReceivedCounter := metrics.GetOrRegisterCounter("counter:action:bytes:received", registry)
		bytesReceivedCounter.Inc(bytesReceivedValue)

		bytesReceived := metrics.GetOrRegisterHistogram("histogram:action:bytes:received", registry, metrics.NewUniformSample(100))
		bytesReceived.Update(int64(result["action:bytes:received"].(int)))
	}

	if result["assertions"] != nil {
		assertionResults := result["assertions"].([]core.AssertionResult)
		for _, result := range assertionResults {
			totalAssertionsCounter := metrics.GetOrRegisterCounter("assertions:total", registry)
			totalAssertionsCounter.Inc(1)
			if result["result"].(bool) == false {
				totalAssertionsFailedCounter := metrics.GetOrRegisterCounter("assertions:failed", registry)
				totalAssertionsFailedCounter.Inc(1)
			}
		}
	}
}
