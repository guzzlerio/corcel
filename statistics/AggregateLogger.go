package statistics

import "github.com/rcrowley/go-metrics"

type AggregateLogger struct {
	LogCounter   func(name string, value int64)
	LogGuage     func(name string, value float64)
	LogTimer     func(name string, value metrics.Timer)
	LogHistogram func(name string, value metrics.Histogram)
	LogMeter     func(name string, value metrics.Meter)
}
