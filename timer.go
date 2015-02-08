package quant

import (
	"sync"
	"time"
)

// TimeUnit represents an enumeration of available
// time units the timer metric supports.
type TimeUnit time.Duration

// All available time units that can be used with a timer metric.
const (
	Nanoseconds  TimeUnit = TimeUnit(time.Nanosecond)
	Microseconds TimeUnit = TimeUnit(time.Microsecond)
	Milliseconds TimeUnit = TimeUnit(time.Millisecond)
	Seconds      TimeUnit = TimeUnit(time.Second)
)

// String returns a string representation of the time unit.
func (tu TimeUnit) String() string {
	switch tu {
	case Nanoseconds:
		return "ns"
	case Microseconds:
		return "µs"
	case Milliseconds:
		return "ms"
	case Seconds:
		return "s"
	default:
		return "*" + time.Duration(tu).String()
	}
}

// Timer represents a time metric which can be used to stop
// duration of tasks. Starting a timer creates a Stopwatch
// for time measurement. If is safe to use a counter concurrently.
type Timer struct {
	metric
	timeUnit TimeUnit
	mtx      sync.Mutex
	snap     *TimerSnapshot
}

func newTimer(name string, unit TimeUnit) *Timer {
	return &Timer{
		metric:   metric{name, unit.String()},
		timeUnit: unit,
		snap:     newTimerSnaphot(name, unit.String()),
	}
}

// Start starts the timer and returns a Stopwatch to measure the duration
// of a specific task.
func (t *Timer) Start() *Stopwatch {
	return newStopwatch(t)
}

func (t *Timer) record(d time.Duration) {
	t.mtx.Lock()
	t.snap.add(float64(d) / float64(t.timeUnit))
	t.mtx.Unlock()
}

func (t *Timer) snapshot() *TimerSnapshot {
	t.mtx.Lock()
	snap := t.snap
	t.snap = newTimerSnaphot(snap.name, snap.unit)
	t.mtx.Unlock()
	return snap
}

// Stopwatch can be used to measure the duration of a specific event
// and report it to the underlying Timer. It is NOT safe to use a
// stopwatch concurrently.
type Stopwatch struct {
	timer *Timer
	t     time.Time
}

func newStopwatch(timer *Timer) *Stopwatch {
	return &Stopwatch{
		timer: timer,
		t:     time.Now(),
	}
}

// Elapsed returns the currently elapsed time since the stopwatch
// was started.
func (sw *Stopwatch) Elapsed() time.Duration {
	return time.Now().Sub(sw.t)
}

// Reset restarts the stopwatch.
func (sw *Stopwatch) Reset() {
	sw.t = time.Now()
}

// Record retrieves the currently elapsed time from Elapsed and report
// this value to the underlying timer. Furthermore the measured duration
// is returned.
func (sw *Stopwatch) Record() time.Duration {
	d := sw.Elapsed()
	sw.timer.record(d)
	return d
}

// TimerSnapshot represents a snapshot of a Timer metric.
// This snapshot type is used during the reporting process.
type TimerSnapshot struct {
	reservoirSnapshot
}

func newTimerSnaphot(name, unit string) *TimerSnapshot {
	return &TimerSnapshot{
		reservoirSnapshot: *newReservoirSnaphot(name, unit),
	}
}
