package statistics

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type Aggregator struct {
	times      []int64
	counters   map[string][]int64
	guages     map[string][]float64
	histograms map[string]map[string][]float64
	meters     map[string]map[string][]float64
	timers     map[string]map[string][]float64
	logger     AggregateLogger
	ticker     *time.Ticker
	registry   metrics.Registry
}

type AggregatorSnapShot struct {
	Times      []int64
	Counters   map[string][]int64
	Guages     map[string][]float64
	Histograms map[string]map[string][]float64
	Meters     map[string]map[string][]float64
	Timers     map[string]map[string][]float64
}

func NewAggregator(registry metrics.Registry) *Aggregator {
	agg := &Aggregator{
		times:      []int64{},
		counters:   map[string][]int64{},
		guages:     map[string][]float64{},
		histograms: map[string]map[string][]float64{},
		meters:     map[string]map[string][]float64{},
		timers:     map[string]map[string][]float64{},
		ticker:     time.NewTicker(time.Second * 2),
		registry:   registry,
	}

	agg.Initialize()

	return agg
}

func (instance *Aggregator) Initialize() {
	instance.logger = AggregateLogger{}
	instance.logger.LogCounter = instance.logCounter
	instance.logger.LogGuage = instance.logGuage
	instance.logger.LogHistogram = instance.logHistogram
	instance.logger.LogMeter = instance.logMeter
	instance.logger.LogTimer = instance.logTimer
}

func (instance *Aggregator) Snapshot() AggregatorSnapShot {
	return AggregatorSnapShot{
		Times:      instance.times,
		Counters:   instance.counters,
		Guages:     instance.guages,
		Histograms: instance.histograms,
		Meters:     instance.meters,
		Timers:     instance.timers,
	}
}

func (instance *Aggregator) logCounter(name string, value int64) {
	if _, ok := instance.counters[name]; !ok {
		instance.counters[name] = []int64{}
	}
	instance.counters[name] = append(instance.counters[name], value)
}

func (instance *Aggregator) logGuage(name string, value float64) {
	if _, ok := instance.guages[name]; !ok {
		instance.guages[name] = []float64{}
	}
	instance.guages[name] = append(instance.guages[name], value)
}

func (instance *Aggregator) logHistogram(name string, value metrics.Histogram) {
	ps := value.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	if _, ok := instance.histograms[name]; !ok {
		instance.histograms[name] = map[string][]float64{}
		instance.histograms[name]["count"] = []float64{}
		instance.histograms[name]["min"] = []float64{}
		instance.histograms[name]["max"] = []float64{}
		instance.histograms[name]["mean"] = []float64{}
		instance.histograms[name]["stddev"] = []float64{}
		instance.histograms[name]["median"] = []float64{}
		instance.histograms[name]["75p"] = []float64{}
		instance.histograms[name]["95p"] = []float64{}
		instance.histograms[name]["99p"] = []float64{}
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
		instance.meters[name]["count"] = []float64{}
		instance.meters[name]["rate1"] = []float64{}
		instance.meters[name]["rate5"] = []float64{}
		instance.meters[name]["rate15"] = []float64{}
		instance.meters[name]["rateMean"] = []float64{}
	}

	instance.meters[name]["count"] = append(instance.meters[name]["count"], float64(value.Count()))
	instance.meters[name]["rate1"] = append(instance.meters[name]["count"], float64(value.Rate1()))
	instance.meters[name]["rate5"] = append(instance.meters[name]["count"], float64(value.Rate5()))
	instance.meters[name]["rate15"] = append(instance.meters[name]["count"], float64(value.Rate15()))
	instance.meters[name]["rateMean"] = append(instance.meters[name]["count"], float64(value.RateMean()))
}

func (instance *Aggregator) logTimer(name string, value metrics.Timer) {
	ps := value.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	if _, ok := instance.meters[name]; !ok {
		instance.meters[name] = map[string][]float64{}
		instance.timers[name]["count"] = []float64{}
		instance.timers[name]["min"] = []float64{}
		instance.timers[name]["max"] = []float64{}
		instance.timers[name]["mean"] = []float64{}
		instance.timers[name]["stddev"] = []float64{}
		instance.timers[name]["median"] = []float64{}
		instance.timers[name]["75p"] = []float64{}
		instance.timers[name]["95p"] = []float64{}
		instance.timers[name]["99p"] = []float64{}
		instance.timers[name]["count"] = []float64{}
		instance.timers[name]["rate1"] = []float64{}
		instance.timers[name]["rate5"] = []float64{}
		instance.timers[name]["rate15"] = []float64{}
		instance.timers[name]["rateMean"] = []float64{}
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
	instance.timers[name]["count"] = append(instance.timers[name]["count"], float64(value.Count()))
	instance.timers[name]["rate1"] = append(instance.timers[name]["count"], float64(value.Rate1()))
	instance.timers[name]["rate5"] = append(instance.timers[name]["count"], float64(value.Rate5()))
	instance.timers[name]["rate15"] = append(instance.timers[name]["count"], float64(value.Rate15()))
	instance.timers[name]["rateMean"] = append(instance.timers[name]["count"], float64(value.RateMean()))
}

func (instance *Aggregator) createSnapshot() {
	instance.times = append(instance.times, time.Now().Unix())
	instance.registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			instance.logCounter(name, metric.Count())
		case metrics.Gauge:
			instance.logGuage(name, float64(metric.Value()))
		case metrics.GaugeFloat64:
			instance.logGuage(name, metric.Value())
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

func (instance *Aggregator) Start() {
	go func() {
		for _ = range instance.ticker.C {
			instance.createSnapshot()
		}
	}()
}

func (instance *Aggregator) Stop() {
	instance.ticker.Stop()
	instance.createSnapshot()
}
