package statistics

import (
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
)

//Aggregator ...
type Aggregator struct {
	times      []int64
	counters   map[string][]int64
	gauges     map[string][]float64
	histograms map[string]map[string][]float64
	meters     map[string]map[string][]float64
	timers     map[string]map[string][]float64
	logger     AggregateLogger
	ticker     *time.Ticker
	registry   metrics.Registry
	mutex      *sync.Mutex
}

//AggregatorSnapShot ...
type AggregatorSnapShot struct {
	Times      []int64
	Counters   map[string][]int64
	Guages     map[string][]float64
	Histograms map[string]map[string][]float64
	Meters     map[string]map[string][]float64
	Timers     map[string]map[string][]float64
}

//NewAggregatorSnapShot ...
func NewAggregatorSnapShot() *AggregatorSnapShot {
	return &AggregatorSnapShot{
		Times:      []int64{},
		Counters:   map[string][]int64{},
		Guages:     map[string][]float64{},
		Histograms: map[string]map[string][]float64{},
		Meters:     map[string]map[string][]float64{},
		Timers:     map[string]map[string][]float64{},
	}
}

//UpdateCounter ...
func (instance *AggregatorSnapShot) UpdateCounter(key string, value int64) {
	if _, ok := instance.Counters[key]; !ok {
		instance.Counters[key] = make([]int64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Counters[key][i] = int64(0)
		}
	}
	instance.Counters[key] = append(instance.Counters[key], value)
}

//UpdateCounters ...
func (instance *AggregatorSnapShot) UpdateCounters(output AggregatorSnapShot) {
	for key, value := range output.Counters {
		instance.UpdateCounter(key, value[len(value)-1])
	}
}

//UpdateGuage ...
func (instance *AggregatorSnapShot) UpdateGuage(key string, value float64) {
	if _, ok := instance.Guages[key]; !ok {
		instance.Guages[key] = make([]float64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Guages[key][i] = float64(0)
		}
	}
	instance.Guages[key] = append(instance.Guages[key], value)
}

//UpdateGuages ...
func (instance *AggregatorSnapShot) UpdateGuages(output AggregatorSnapShot) {
	for key, value := range output.Guages {
		instance.UpdateGuage(key, value[len(value)-1])
	}
}

//UpdateHistogram ...
func (instance *AggregatorSnapShot) UpdateHistogram(key string, subKey string, value float64) {
	if _, ok := instance.Histograms[key]; !ok {
		instance.Histograms[key] = map[string][]float64{}
	}
	if _, ok := instance.Histograms[key][subKey]; !ok {
		instance.Histograms[key][subKey] = make([]float64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Histograms[key][subKey][i] = float64(0)
		}
	}
	instance.Histograms[key][subKey] = append(instance.Histograms[key][subKey], value)
}

//UpdateHistograms ...
func (instance *AggregatorSnapShot) UpdateHistograms(output AggregatorSnapShot) {
	for key, value := range output.Histograms {
		for subKey, subValue := range value {
			instance.UpdateHistogram(key, subKey, subValue[len(subValue)-1])
		}
	}
}

//UpdateMeter ...
func (instance *AggregatorSnapShot) UpdateMeter(key string, subKey string, value float64) {
	if _, ok := instance.Meters[key]; !ok {
		instance.Meters[key] = map[string][]float64{}
	}
	if _, ok := instance.Meters[key][subKey]; !ok {
		instance.Meters[key][subKey] = make([]float64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Meters[key][subKey][i] = float64(0)
		}
	}
	instance.Meters[key][subKey] = append(instance.Meters[key][subKey], value)
}

//UpdateMeters ...
func (instance *AggregatorSnapShot) UpdateMeters(output AggregatorSnapShot) {
	for key, value := range output.Meters {
		for subKey, subValue := range value {
			instance.UpdateMeter(key, subKey, subValue[len(subValue)-1])
		}
	}
}

//UpdateTimer ...
func (instance *AggregatorSnapShot) UpdateTimer(key string, subKey string, value float64) {
	if _, ok := instance.Timers[key]; !ok {
		instance.Timers[key] = map[string][]float64{}
	}
	if _, ok := instance.Timers[key][subKey]; !ok {
		instance.Timers[key][subKey] = make([]float64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Timers[key][subKey][i] = float64(0)
		}
	}
	instance.Timers[key][subKey] = append(instance.Timers[key][subKey], value)
}

//UpdateTimers ...
func (instance *AggregatorSnapShot) UpdateTimers(output AggregatorSnapShot) {
	for key, value := range output.Timers {
		for subKey, subValue := range value {
			instance.UpdateTimer(key, subKey, subValue[len(subValue)-1])
		}
	}
}

//UpdateTime ...
func (instance *AggregatorSnapShot) UpdateTime(value int64) {
	instance.Times = append(instance.Times, value)
}

//ExecutionSummary ...
type ExecutionSummary struct {
	TotalRequests          float64
	TotalErrors            float64
	Availability           float64
	RunningTime            string
	Throughput             float64
	TotalBytesSent         int64
	TotalBytesReceived     int64
	MeanResponseTime       float64
	MinResponseTime        float64
	MaxResponseTime        float64
	TotalAssertions        int64
	TotalAssertionFailures int64
}

//CreateSummary ...
func CreateSummary(snapshot AggregatorSnapShot) ExecutionSummary {

	lastTime := time.Unix(0, snapshot.Times[len(snapshot.Times)-1])
	firstTime := time.Unix(0, snapshot.Times[0])
	duration := lastTime.Sub(firstTime)

	counts := snapshot.Meters["action:throughput"]["count"]
	count := counts[len(counts)-1]

	errors := snapshot.Meters["action:error"]["count"]
	errorCount := errors[len(errors)-1]

	rates := snapshot.Meters["action:throughput"]["rateMean"]
	rate := rates[len(rates)-1]

	var availability float64
	if errorCount > 0 {
		availability = (1 - (float64(errorCount) / float64(count))) * 100
	} else {
		availability = 100
	}

	var bytesSentCount = int64(0)
	var bytesReceivedCount = int64(0)
	var totalAssertionsCount = int64(0)
	var totalAssertionFailuresCount = int64(0)

	bytesSent := snapshot.Counters["action:bytes:sent"]

	if bytesSent != nil {
		bytesSentCount = bytesSent[len(bytesSent)-1]
	}

	bytesReceived := snapshot.Counters["action:bytes:received"]
	if bytesReceived != nil {
		bytesReceivedCount = bytesReceived[len(bytesReceived)-1]
	}

	totalAssertions := snapshot.Counters["assertions:total"]
	if totalAssertions != nil {
		totalAssertionsCount = totalAssertions[len(totalAssertions)-1]
	}

	totalAssertionsFailed := snapshot.Counters["assertions:failed"]
	if totalAssertionsFailed != nil {
		totalAssertionFailuresCount = totalAssertionsFailed[len(totalAssertionsFailed)-1]
	}

	responseMeanTimes := snapshot.Timers["action:duration"]["mean"]
	responseMeanTime := responseMeanTimes[len(responseMeanTimes)-1]

	responseMinTimes := snapshot.Timers["action:duration"]["min"]
	responseMinTime := responseMinTimes[len(responseMinTimes)-1]

	responseMaxTimes := snapshot.Timers["action:duration"]["max"]
	responseMaxTime := responseMaxTimes[len(responseMaxTimes)-1]

	return ExecutionSummary{
		RunningTime:            duration.String(),
		TotalRequests:          count,
		TotalErrors:            errorCount,
		Availability:           availability,
		Throughput:             rate,
		TotalBytesSent:         bytesSentCount,
		TotalBytesReceived:     bytesReceivedCount,
		MeanResponseTime:       responseMeanTime,
		MinResponseTime:        responseMinTime,
		MaxResponseTime:        responseMaxTime,
		TotalAssertions:        totalAssertionsCount,
		TotalAssertionFailures: totalAssertionFailuresCount,
	}
}

//NewAggregator ...
func NewAggregator(registry metrics.Registry) *Aggregator {
	agg := &Aggregator{
		times:      []int64{},
		counters:   map[string][]int64{},
		gauges:     map[string][]float64{},
		histograms: map[string]map[string][]float64{},
		meters:     map[string]map[string][]float64{},
		timers:     map[string]map[string][]float64{},
		ticker:     time.NewTicker(time.Second * 2),
		registry:   registry,
		mutex:      &sync.Mutex{},
	}

	agg.Initialize()

	return agg
}

//Initialize ...
func (instance *Aggregator) Initialize() {
	instance.logger = AggregateLogger{}
	instance.logger.LogCounter = instance.logCounter
	instance.logger.LogGuage = instance.logGuage
	instance.logger.LogHistogram = instance.logHistogram
	instance.logger.LogMeter = instance.logMeter
	instance.logger.LogTimer = instance.logTimer
}

//Snapshot ...
func (instance *Aggregator) Snapshot() AggregatorSnapShot {
	return AggregatorSnapShot{
		Times:      instance.times,
		Counters:   instance.counters,
		Guages:     instance.gauges,
		Histograms: instance.histograms,
		Meters:     instance.meters,
		Timers:     instance.timers,
	}
}

func (instance *Aggregator) logCounter(name string, value int64) {
	name = strings.Replace(name, "counter:", "", -1)
	if _, ok := instance.counters[name]; !ok {
		instance.counters[name] = make([]int64, len(instance.times)-1)
		for i := 0; i < len(instance.times)-1; i++ {
			instance.counters[name][i] = int64(0)
		}
	}
	instance.counters[name] = append(instance.counters[name], value)
}

func (instance *Aggregator) logGuage(name string, value float64) {
	if _, ok := instance.gauges[name]; !ok {
		instance.gauges[name] = make([]float64, len(instance.times)-1)
		for i := 0; i < len(instance.times)-1; i++ {
			instance.gauges[name][i] = float64(0)
		}
	}
	instance.gauges[name] = append(instance.gauges[name], value)
}

func (instance *Aggregator) logHistogram(name string, value metrics.Histogram) {
	name = strings.Replace(name, "histogram:", "", -1)
	ps := value.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	if _, ok := instance.histograms[name]; !ok {
		instance.histograms[name] = map[string][]float64{}
		instance.histograms[name]["count"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["min"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["max"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["mean"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["stddev"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["median"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["75p"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["95p"] = make([]float64, len(instance.times)-1)
		instance.histograms[name]["99p"] = make([]float64, len(instance.times)-1)

		for i := 0; i < len(instance.times)-1; i++ {
			instance.histograms[name]["count"][i] = float64(0)
			instance.histograms[name]["min"][i] = float64(0)
			instance.histograms[name]["max"][i] = float64(0)
			instance.histograms[name]["mean"][i] = float64(0)
			instance.histograms[name]["stddev"][i] = float64(0)
			instance.histograms[name]["median"][i] = float64(0)
			instance.histograms[name]["75p"][i] = float64(0)
			instance.histograms[name]["95p"][i] = float64(0)
			instance.histograms[name]["99p"][i] = float64(0)
		}
	}
	instance.histograms[name]["count"] = append(instance.histograms[name]["count"], float64(value.Count()))
	instance.histograms[name]["min"] = append(instance.histograms[name]["min"], float64(value.Min()))
	instance.histograms[name]["max"] = append(instance.histograms[name]["max"], float64(value.Max()))
	instance.histograms[name]["mean"] = append(instance.histograms[name]["mean"], value.Mean())
	instance.histograms[name]["stddev"] = append(instance.histograms[name]["stddev"], value.StdDev())
	instance.histograms[name]["median"] = append(instance.histograms[name]["median"], ps[0])
	instance.histograms[name]["75p"] = append(instance.histograms[name]["75p"], ps[1])
	instance.histograms[name]["95p"] = append(instance.histograms[name]["95p"], ps[2])
	instance.histograms[name]["99p"] = append(instance.histograms[name]["99p"], ps[3])
}

func (instance *Aggregator) logMeter(name string, value metrics.Meter) {
	if _, ok := instance.meters[name]; !ok {
		instance.meters[name] = map[string][]float64{}
		instance.meters[name]["count"] = make([]float64, len(instance.times)-1)
		instance.meters[name]["rate1"] = make([]float64, len(instance.times)-1)
		instance.meters[name]["rate5"] = make([]float64, len(instance.times)-1)
		instance.meters[name]["rate15"] = make([]float64, len(instance.times)-1)
		instance.meters[name]["rateMean"] = make([]float64, len(instance.times)-1)

		for i := 0; i < len(instance.times)-1; i++ {
			instance.meters[name]["count"][i] = float64(0)
			instance.meters[name]["rate1"][i] = float64(0)
			instance.meters[name]["rate5"][i] = float64(0)
			instance.meters[name]["rate15"][i] = float64(0)
			instance.meters[name]["rateMean"][i] = float64(0)
		}
	}

	instance.meters[name]["count"] = append(instance.meters[name]["count"], float64(value.Count()))
	instance.meters[name]["rate1"] = append(instance.meters[name]["rate1"], float64(value.Rate1()))
	instance.meters[name]["rate5"] = append(instance.meters[name]["rate5"], float64(value.Rate5()))
	instance.meters[name]["rate15"] = append(instance.meters[name]["rate15"], float64(value.Rate15()))
	instance.meters[name]["rateMean"] = append(instance.meters[name]["rateMean"], float64(value.RateMean()))
}

func (instance *Aggregator) logTimer(name string, value metrics.Timer) {
	ps := value.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	if _, ok := instance.timers[name]; !ok {
		instance.timers[name] = map[string][]float64{}
		instance.timers[name]["count"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["min"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["max"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["mean"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["stddev"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["median"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["75p"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["95p"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["99p"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["count"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["rate1"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["rate5"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["rate15"] = make([]float64, len(instance.times)-1)
		instance.timers[name]["rateMean"] = make([]float64, len(instance.times)-1)

		for i := 0; i < len(instance.times)-1; i++ {
			instance.timers[name]["count"][i] = float64(0)
			instance.timers[name]["min"][i] = float64(0)
			instance.timers[name]["max"][i] = float64(0)
			instance.timers[name]["mean"][i] = float64(0)
			instance.timers[name]["stddev"][i] = float64(0)
			instance.timers[name]["median"][i] = float64(0)
			instance.timers[name]["75p"][i] = float64(0)
			instance.timers[name]["95p"][i] = float64(0)
			instance.timers[name]["99p"][i] = float64(0)
			instance.timers[name]["count"][i] = float64(0)
			instance.timers[name]["rate1"][i] = float64(0)
			instance.timers[name]["rate5"][i] = float64(0)
			instance.timers[name]["rate15"][i] = float64(0)
			instance.timers[name]["rateMean"][i] = float64(0)
		}
	}

	instance.timers[name]["count"] = append(instance.timers[name]["count"], float64(value.Count()))
	instance.timers[name]["min"] = append(instance.timers[name]["min"], float64(value.Min()))
	instance.timers[name]["max"] = append(instance.timers[name]["max"], float64(value.Max()))
	instance.timers[name]["mean"] = append(instance.timers[name]["mean"], value.Mean())
	instance.timers[name]["stddev"] = append(instance.timers[name]["stddev"], value.StdDev())
	instance.timers[name]["median"] = append(instance.timers[name]["median"], ps[0])
	instance.timers[name]["75p"] = append(instance.timers[name]["75p"], ps[1])
	instance.timers[name]["95p"] = append(instance.timers[name]["95p"], ps[2])
	instance.timers[name]["99p"] = append(instance.timers[name]["99p"], ps[3])
	instance.timers[name]["rate1"] = append(instance.timers[name]["rate1"], float64(value.Rate1()))
	instance.timers[name]["rate5"] = append(instance.timers[name]["rate5"], float64(value.Rate5()))
	instance.timers[name]["rate15"] = append(instance.timers[name]["rate15"], float64(value.Rate15()))
	instance.timers[name]["rateMean"] = append(instance.timers[name]["rateMean"], float64(value.RateMean()))
}

func (instance *Aggregator) createSnapshot() {
	timeToLog := time.Now().UnixNano()
	if len(instance.times) > 1 && instance.times[len(instance.times)-1] == int64(timeToLog) {
		return
	}

	instance.mutex.Lock()
	instance.times = append(instance.times, timeToLog)
	instance.registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			counter := metric.Snapshot()
			instance.logCounter(name, counter.Count())
		case metrics.Gauge:
			instance.logGuage(name, float64(metric.Snapshot().Value()))
		case metrics.GaugeFloat64:
			instance.logGuage(name, metric.Snapshot().Value())
		case metrics.Histogram:
			h := metric.Snapshot()
			instance.logHistogram(name, h)
		case metrics.Meter:
			m := metric.Snapshot()
			instance.logMeter(name, m)
		case metrics.Timer:
			t := metric.Snapshot()
			instance.logTimer(name, t)
		}
	})
	instance.mutex.Unlock()
}

//Start ...
func (instance *Aggregator) Start() {
	instance.createSnapshot()
	go func() {
		for _ = range instance.ticker.C {
			instance.createSnapshot()
		}
	}()
}

//Stop ...
func (instance *Aggregator) Stop() {
	instance.createSnapshot()
	instance.ticker.Stop()
}
