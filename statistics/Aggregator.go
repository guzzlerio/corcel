package statistics

import (
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
)

//Aggregator ...
type Aggregator struct {
	times      []int64
	counters   map[string][]int64
	gauges     map[string][]float64
	histograms map[string]map[string][]int64
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
	Gauges     map[string][]float64
	Histograms map[string]map[string][]int64
	Meters     map[string]map[string][]float64
	Timers     map[string]map[string][]float64
}

//NewAggregatorSnapShot ...
func NewAggregatorSnapShot() *AggregatorSnapShot {
	return &AggregatorSnapShot{
		Times:      []int64{},
		Counters:   map[string][]int64{},
		Gauges:     map[string][]float64{},
		Histograms: map[string]map[string][]int64{},
		Meters:     map[string]map[string][]float64{},
		Timers:     map[string]map[string][]float64{},
	}
}

func (instance *AggregatorSnapShot) updateCounter(key string, value int64) {
	if _, ok := instance.Counters[key]; !ok {
		instance.Counters[key] = make([]int64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Counters[key][i] = int64(0)
		}
	}
	instance.Counters[key] = append(instance.Counters[key], value)
}

func (instance *AggregatorSnapShot) updateCounters(output AggregatorSnapShot) {
	for key, value := range output.Counters {
		instance.updateCounter(key, value[len(value)-1])
	}
}

func (instance *AggregatorSnapShot) updateGauge(key string, value float64) {
	if _, ok := instance.Gauges[key]; !ok {
		instance.Gauges[key] = make([]float64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Gauges[key][i] = float64(0)
		}
	}
	instance.Gauges[key] = append(instance.Gauges[key], value)
}

func (instance *AggregatorSnapShot) updateGauges(output AggregatorSnapShot) {
	for key, value := range output.Gauges {
		instance.updateGauge(key, value[len(value)-1])
	}
}

func (instance *AggregatorSnapShot) updateHistogram(key string, subKey string, value int64) {
	if _, ok := instance.Histograms[key]; !ok {
		instance.Histograms[key] = map[string][]int64{}
	}
	if _, ok := instance.Histograms[key][subKey]; !ok {
		instance.Histograms[key][subKey] = make([]int64, len(instance.Times))
		for i := 0; i < len(instance.Times)-1; i++ {
			instance.Histograms[key][subKey][i] = int64(0)
		}
	}
	instance.Histograms[key][subKey] = append(instance.Histograms[key][subKey], value)
}

func (instance *AggregatorSnapShot) updateHistograms(output AggregatorSnapShot) {
	for key, value := range output.Histograms {
		for subKey, subValue := range value {
			instance.updateHistogram(key, subKey, subValue[len(subValue)-1])
		}
	}
}

func (instance *AggregatorSnapShot) updateMeter(key string, subKey string, value float64) {
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

func (instance *AggregatorSnapShot) updateMeters(output AggregatorSnapShot) {
	for key, value := range output.Meters {
		for subKey, subValue := range value {
			instance.updateMeter(key, subKey, subValue[len(subValue)-1])
		}
	}
}

func (instance *AggregatorSnapShot) updateTimer(key string, subKey string, value float64) {
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

func (instance *AggregatorSnapShot) updateTimers(output AggregatorSnapShot) {
	for key, value := range output.Timers {
		for subKey, subValue := range value {
			instance.updateTimer(key, subKey, subValue[len(subValue)-1])
		}
	}
}

func (instance *AggregatorSnapShot) updateTime(value int64) {
	instance.Times = append(instance.Times, value)
}

//Update ...
func (instance *AggregatorSnapShot) Update(output AggregatorSnapShot) {
	instance.updateCounters(output)
	instance.updateGauges(output)
	instance.updateHistograms(output)
	instance.updateMeters(output)
	instance.updateTimers(output)
	instance.updateTime(output.Times[len(output.Times)-1])
}

//IncrementCounter ...
func IncrementCounter(registry metrics.Registry, key string, value int64) {
	counter := metrics.GetOrRegisterCounter(key, registry)
	counter.Inc(value)
}

//NewAggregator ...
func NewAggregator(registry metrics.Registry) *Aggregator {
	agg := &Aggregator{
		times:      []int64{},
		counters:   map[string][]int64{},
		gauges:     map[string][]float64{},
		histograms: map[string]map[string][]int64{},
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
		Gauges:     instance.gauges,
		Histograms: instance.histograms,
		Meters:     instance.meters,
		Timers:     instance.timers,
	}
}

func (instance *Aggregator) logCounter(name string, value int64) {
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
	ps := value.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	if _, ok := instance.histograms[name]; !ok {
		instance.histograms[name] = map[string][]int64{}
		instance.histograms[name]["count"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["min"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["max"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["mean"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["stddev"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["median"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["75p"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["95p"] = make([]int64, len(instance.times)-1)
		instance.histograms[name]["99p"] = make([]int64, len(instance.times)-1)

		for i := 0; i < len(instance.times)-1; i++ {
			instance.histograms[name]["count"][i] = int64(0)
			instance.histograms[name]["min"][i] = int64(0)
			instance.histograms[name]["max"][i] = int64(0)
			instance.histograms[name]["mean"][i] = int64(0)
			instance.histograms[name]["stddev"][i] = int64(0)
			instance.histograms[name]["median"][i] = int64(0)
			instance.histograms[name]["75p"][i] = int64(0)
			instance.histograms[name]["95p"][i] = int64(0)
			instance.histograms[name]["99p"][i] = int64(0)
		}
	}
	instance.histograms[name]["count"] = append(instance.histograms[name]["count"], int64(value.Count()))
	instance.histograms[name]["min"] = append(instance.histograms[name]["min"], int64(value.Min()))
	instance.histograms[name]["max"] = append(instance.histograms[name]["max"], int64(value.Max()))
	instance.histograms[name]["mean"] = append(instance.histograms[name]["mean"], int64(value.Mean()))
	instance.histograms[name]["stddev"] = append(instance.histograms[name]["stddev"], int64(value.StdDev()))
	instance.histograms[name]["median"] = append(instance.histograms[name]["median"], int64(ps[0]))
	instance.histograms[name]["75p"] = append(instance.histograms[name]["75p"], int64(ps[1]))
	instance.histograms[name]["95p"] = append(instance.histograms[name]["95p"], int64(ps[2]))
	instance.histograms[name]["99p"] = append(instance.histograms[name]["99p"], int64(ps[3]))
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

	instance.times = append(instance.times, timeToLog)
	instance.registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			counter := metric.Snapshot()
			instance.logCounter(name, counter.Count())
			metric.Clear()
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
