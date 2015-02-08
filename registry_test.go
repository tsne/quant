package quant

import (
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	reg := NewRegistry("reg")

	c := reg.NewCounter("my-counter")
	g := reg.NewGauge("my-gauge", func() float64 { return 0 })
	tm := reg.NewTimer("my-timer", Milliseconds)

	switch {
	case c == nil:
		t.Error("no timer in registry")
	case c != reg.Counter("my-counter"):
		t.Error("wrong timer in registry")
	}

	switch {
	case g == nil:
		t.Error("no gauge in registry")
	case g != reg.Gauge("my-gauge"):
		t.Error("wrong gauge in registry")
	}

	switch {
	case tm == nil:
		t.Error("no timer in registry")
	case tm != reg.Timer("my-timer"):
		t.Error("wrong timer in registry")
	}
}

func TestRegistryExistingMetric(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("registry does not panic for existing metrics")
		}
	}()

	reg := NewRegistry("reg")
	reg.NewCounter("m")
	reg.NewGauge("m", func() float64 { return 0 })
}

type testReporter struct {
	reportCounters func(string, []*CounterSnapshot) error
	reportGauges   func(string, []*GaugeSnapshot) error
	reportTimers   func(string, []*TimerSnapshot) error
}

func (r *testReporter) ReportCounters(registryName string, counters []*CounterSnapshot) error {
	return r.reportCounters(registryName, counters)
}

func (r *testReporter) ReportGauges(registryName string, gauges []*GaugeSnapshot) error {
	return r.reportGauges(registryName, gauges)
}

func (r *testReporter) ReportTimers(registryName string, timers []*TimerSnapshot) error {
	return r.reportTimers(registryName, timers)
}

func TestRegistryReporting(t *testing.T) {
	reg := NewRegistry("reg")

	reg.NewCounter("my-counter")
	reg.Counter("my-counter").Increment()
	reg.Counter("my-counter").Increment()

	reg.NewGauge("my-gauge", func() float64 { return 7.0 })

	reg.NewTimer("my-timer", Milliseconds)
	sw := reg.Timer("my-timer").Start()
	time.Sleep(10 * time.Millisecond)
	sw.Record()

	reg.Report(&testReporter{
		reportCounters: func(registryName string, counters []*CounterSnapshot) error {
			if len(counters) != 1 {
				t.Errorf("wrong number of counters: %d (1 expected)", len(counters))
			}
			if counters[0].Value() != 2 {
				t.Errorf("wrong counter value: %d (2 expected)", counters[0].Value())
			}
			return nil
		},
		reportGauges: func(registryName string, gauges []*GaugeSnapshot) error {
			if len(gauges) != 1 {
				t.Errorf("wrong number of gauges: %d (1 expected)", len(gauges))
			}
			if gauges[0].Value() != 7.0 {
				t.Errorf("wrong counter value: %f (7 expected)", gauges[0].Value())
			}
			return nil
		},
		reportTimers: func(registryName string, timers []*TimerSnapshot) error {
			if len(timers) != 1 {
				t.Errorf("wrong number of timers: %d (1 expected)", len(timers))
			}
			if timers[0].Count() != 1 {
				t.Errorf("wrong number of times: %d (1 expected)", timers[0].Count())
			}
			if timers[0].Minimum() < 10 || timers[0].Maximum() < 10 || timers[0].Average() < 10 {
				t.Errorf("wrong timer value: %f (at least 10 expected)", timers[0].Minimum())
			}
			if timers[0].Variance() != 0 {
				t.Errorf("wrong timer variance: %f (0 expected)", timers[0].Variance())
			}
			return nil
		},
	})
}
