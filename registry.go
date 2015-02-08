package quant

import (
	"fmt"
	"sync"
)

// Registry represents a collection of metrics. Each metric
// in this collection must have a unique name which identifies
// the respective metric. A registry has the ability to report
// its metrics to a set of reporters.
// It is safe to use a metrics registry concurrently.
type Registry struct {
	name        string
	mtx         sync.RWMutex
	metricNames map[string]struct{}
	counters    map[string]*Counter
	gauges      map[string]*Gauge
	timers      map[string]*Timer
}

// NewRegistry creates a new registry with the specified name.
// The name act as an identifier during a reporting.
func NewRegistry(name string) *Registry {
	return &Registry{
		name:        name,
		metricNames: make(map[string]struct{}),
		counters:    make(map[string]*Counter),
		gauges:      make(map[string]*Gauge),
		timers:      make(map[string]*Timer),
	}
}

// Name returns the name of the registry.
func (r *Registry) Name() string {
	return r.name
}

// NewCounter adds a new counter metric to the registry.
// If the given name already exists this function will panic.
func (r *Registry) NewCounter(name string) *Counter {
	return r.NewCounterWithUnit(name, "")
}

// NewCounterWithUnit adds a new counter metric with the specified unit
// to the registry.
// If the given name already exists this function will panic.
func (r *Registry) NewCounterWithUnit(name, unit string) *Counter {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, exists := r.metricNames[name]; exists {
		panic(fmt.Errorf("metric already exists: %s", name))
	}

	counter := newCounter(name, unit)
	r.counters[name] = counter
	r.metricNames[name] = struct{}{}
	return counter
}

// Counter retrieves the counter with the given name. If no such
// counter exists nil will be returned.
func (r *Registry) Counter(name string) *Counter {
	r.mtx.RLock()
	counter := r.counters[name]
	r.mtx.RUnlock()
	return counter
}

// NewGauge adds a new gauge metric to the registry.
// If the given name already exists this function will panic.
func (r *Registry) NewGauge(name string, reader GaugeReader) *Gauge {
	return r.NewGaugeWithUnit(name, "", reader)
}

// NewGaugeWithUnit adds a new gauge metric with the specified unit
// to the registry.
// If the given name already exists this function will panic.
func (r *Registry) NewGaugeWithUnit(name, unit string, reader GaugeReader) *Gauge {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, exists := r.metricNames[name]; exists {
		panic(fmt.Errorf("metric already exists: %s", name))
	}

	gauge := newGauge(name, unit, reader)
	r.gauges[name] = gauge
	r.metricNames[name] = struct{}{}
	return gauge
}

// Gauge retrieves the gauge with the given name. If no such
// gauge exists nil will be returned.
func (r *Registry) Gauge(name string) *Gauge {
	r.mtx.RLock()
	gauge := r.gauges[name]
	r.mtx.RUnlock()
	return gauge
}

// NewTimer adds a new timer metric with the specified unit
// to the registry.
// If the given name already exists this function will panic.
func (r *Registry) NewTimer(name string, unit TimeUnit) *Timer {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, exists := r.metricNames[name]; exists {
		panic(fmt.Errorf("metric already exists: %s", name))
	}

	timer := newTimer(name, unit)
	r.timers[name] = timer
	r.metricNames[name] = struct{}{}
	return timer
}

// Timer retrieves the timer with the given name. If no such
// timer exists nil will be returned.
func (r *Registry) Timer(name string) *Timer {
	r.mtx.RLock()
	timer := r.timers[name]
	r.mtx.RUnlock()
	return timer
}

// Contains checks if a given metric name exists in this registry.
func (r *Registry) Contains(name string) bool {
	r.mtx.RLock()
	_, found := r.metricNames[name]
	r.mtx.RUnlock()
	return found
}

// Report writes the snapshots of all registered metrics to the
// given reporters. If one reporter returns an error during execution
// this error will be returned without executing the followwing reporters.
func (r *Registry) Report(reporters ...Reporter) error {
	if len(reporters) == 0 {
		return nil
	}

	r.mtx.RLock()
	counters := r.counterSnapshots()
	gauges := r.gaugeSnapshots()
	timers := r.timerSnapshots()
	r.mtx.RUnlock()

	for _, reporter := range reporters {
		if len(counters) != 0 {
			if err := reporter.ReportCounters(r.name, counters); err != nil {
				return err
			}
		}
		if len(gauges) != 0 {
			if err := reporter.ReportGauges(r.name, gauges); err != nil {
				return err
			}
		}
		if len(timers) != 0 {
			if err := reporter.ReportTimers(r.name, timers); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Registry) counterSnapshots() []*CounterSnapshot {
	snapshots := make([]*CounterSnapshot, len(r.counters))
	idx := 0
	for _, counter := range r.counters {
		snapshots[idx] = counter.snapshot()
		idx++
	}
	return snapshots
}

func (r *Registry) gaugeSnapshots() []*GaugeSnapshot {
	snapshots := make([]*GaugeSnapshot, len(r.gauges))
	idx := 0
	for _, gauge := range r.gauges {
		snapshots[idx] = gauge.snapshot()
		idx++
	}
	return snapshots
}

func (r *Registry) timerSnapshots() []*TimerSnapshot {
	snapshots := make([]*TimerSnapshot, len(r.timers))
	idx := 0
	for _, timer := range r.timers {
		snapshots[idx] = timer.snapshot()
		idx++
	}
	return snapshots
}
