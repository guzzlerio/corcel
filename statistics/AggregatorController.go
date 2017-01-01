package statistics

import (
	"sync"
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

//AggregatorController ...
type AggregatorController struct {
	snapShotCreateRequest chan struct{}
	snapShotQueryRequest  chan struct{}
	snapShotQueryResponse chan AggregatorSnapShot
	stoppingSignal        chan struct{}
	stoppedSignal         chan struct{}
	registry              metrics.Registry
	ticker                *time.Ticker
	final                 *AggregatorSnapShot
}

//CreateAggregatorController ...
func CreateAggregatorController(registry metrics.Registry) *AggregatorController {
	return &AggregatorController{
		snapShotCreateRequest: make(chan struct{}),
		snapShotQueryRequest:  make(chan struct{}),
		snapShotQueryResponse: make(chan AggregatorSnapShot),
		stoppingSignal:        make(chan struct{}),
		stoppedSignal:         make(chan struct{}),
		registry:              registry,
		ticker:                time.NewTicker(time.Second * 2),
	}
}

//Start ...
func (instance *AggregatorController) Start() {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		agg := &Aggregator{
			times:      []int64{},
			counters:   map[string][]int64{},
			gauges:     map[string][]float64{},
			histograms: map[string]map[string][]int64{},
			meters:     map[string]map[string][]float64{},
			timers:     map[string]map[string][]float64{},
			ticker:     time.NewTicker(time.Second * 2),
			registry:   instance.registry,
			mutex:      &sync.Mutex{},
		}

		agg.Initialize()

		var run = true
		wg.Done()
		for run == true {
			select {
			case <-instance.snapShotCreateRequest:
				agg.createSnapshot()
				break
			case <-instance.snapShotQueryRequest:
				var snapshot = AggregatorSnapShot{
					Times:      agg.times,
					Counters:   agg.counters,
					Gauges:     agg.gauges,
					Histograms: agg.histograms,
					Meters:     agg.meters,
					Timers:     agg.timers,
				}
				instance.snapShotQueryResponse <- snapshot
				break
			case <-instance.stoppingSignal:
				run = false
				break
			default:
				time.Sleep(1)
			}
		}
		instance.ticker.Stop()
		close(instance.stoppedSignal)
	}()
	go func() {
		instance.snapShotCreateRequest <- struct{}{}
		wg.Done()
		for _ = range instance.ticker.C {
			instance.snapShotCreateRequest <- struct{}{}
		}
	}()
	wg.Wait()
}

//Stop ...
func (instance *AggregatorController) Stop() {
	instance.snapShotCreateRequest <- struct{}{}
	instance.snapShotQueryRequest <- struct{}{}
	var snapshot = <-instance.snapShotQueryResponse
	instance.final = &snapshot
	instance.stoppingSignal <- struct{}{}
	<-instance.stoppedSignal
}

//Initialize ...
func (instance *AggregatorController) Initialize() {

}

//Snapshot ...
func (instance *AggregatorController) Snapshot() AggregatorSnapShot {
	if instance.final != nil {
		return *instance.final
	}
	instance.snapShotQueryRequest <- struct{}{}
	return <-instance.snapShotQueryResponse
}
