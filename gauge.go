package quant

import (
	"sync"
)

// GaugeReader represents a function that returns the gauge
// value whenever it is called.
type GaugeReader func() float64

// Gauge represents a float64 metric which value is determined
// by the underlying GaugeReader. Calling the gauge reader is
// protected by a mutex. So it is safe to retrieve the gauge's
// value concurrently.
type Gauge struct {
	metric
	mtx    sync.Mutex
	reader GaugeReader
}

func newGauge(name, unit string, reader GaugeReader) *Gauge {
	return &Gauge{
		metric: metric{name, unit},
		reader: reader,
	}
}

// Value returns the current value of tha underlying GaugeReader.
func (g *Gauge) Value() float64 {
	g.mtx.Lock()
	val := g.reader()
	g.mtx.Unlock()
	return val
}

func (g *Gauge) snapshot() *GaugeSnapshot {
	return &GaugeSnapshot{
		snapshot: snapshot{g.name, g.unit},
		value:    g.Value(),
	}
}

// GaugeSnapshot represents a snapshot of a Gauge metric.
// This snapshot type is used during the reporting process.
type GaugeSnapshot struct {
	snapshot
	value float64
}

// Value returns the snapshot value of the underlying gauge.
func (s *GaugeSnapshot) Value() float64 {
	return s.value
}
