package quant

import (
	"fmt"
)

// Reporter is an interface that is used by a registry to write
// metric snapshots to an specified location.
type Reporter interface {
	ReportCounters(registryName string, counters []*CounterSnapshot) error
	ReportGauges(registryName string, gauges []*GaugeSnapshot) error
	ReportTimers(registryName string, timers []*TimerSnapshot) error
}

// NullReporter is a Reporter implementation that does nothing. Each function
// simply returns a nil as an error.
var NullReporter = nullReporter{}

type nullReporter struct{}

func (r nullReporter) ReportCounters(registryName string, counters []*CounterSnapshot) error {
	return nil
}

func (r nullReporter) ReportGauges(registryName string, gauges []*GaugeSnapshot) error {
	return nil
}

func (r nullReporter) ReportTimers(registryName string, timers []*TimerSnapshot) error {
	return nil
}

// StdoutReporter is a Reporter implementation that simply writes the
// metric snapshots to the standard output.
var StdoutReporter = stdoutReporter{}

type stdoutReporter struct{}

func (r stdoutReporter) ReportCounters(registryName string, counters []*CounterSnapshot) error {
	fmt.Printf("counters of %s\n", registryName)
	for _, c := range counters {
		fmt.Printf("  %s: %d%s\n", c.Name(), c.Value(), c.Unit())
	}
	return nil
}

func (r stdoutReporter) ReportGauges(registryName string, gauges []*GaugeSnapshot) error {
	fmt.Printf("gauges of %s\n", registryName)
	for _, g := range gauges {
		fmt.Printf("  %s: %f%s\n", g.Name(), g.Value(), g.Unit())
	}
	return nil
}

func (r stdoutReporter) ReportTimers(registryName string, timers []*TimerSnapshot) error {
	fmt.Printf("timers of %s\n", registryName)
	for _, t := range timers {
		fmt.Printf("  %s: min=%f%s, max=%f%s, avg=%f%s, dev=%f\n",
			t.Name(), t.Minimum(), t.Unit(), t.Maximum(), t.Unit(), t.Average(), t.Unit(), t.StdDeviation())
	}
	return nil
}
