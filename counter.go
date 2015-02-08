package quant

import (
	"sync/atomic"
)

// Counter represents an int64 metric which can be incremented
// and decremented. It is safe to use a counter concurrently.
type Counter struct {
	metric
	value int64
}

func newCounter(name string, unit string) *Counter {
	return &Counter{
		metric: metric{name, unit},
		value:  0,
	}
}

// Value returns the current int64 value of the counter.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Increment increases the counter by one.
func (c *Counter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement decreases the counter by one.
func (c *Counter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Reset sets the counter back to zero.
func (c *Counter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

func (c *Counter) snapshot() *CounterSnapshot {
	return &CounterSnapshot{
		snapshot: snapshot{c.name, c.unit},
		value:    c.Value(),
	}
}

// CounterSnapshot represents a snapshot of a Counter metric.
// This snapshot type is used during the reporting process.
type CounterSnapshot struct {
	snapshot
	value int64
}

// Value returns the snapshot value of the underlying counter.
func (s *CounterSnapshot) Value() int64 {
	return s.value
}
