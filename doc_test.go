package quant

import (
	"time"
)

func ExampleCounter() {
	registry := NewRegistry("my-registry")
	counter := registry.NewCounter("my-counter")

	// use the counter
	counter.Increment() // or counter.Decrement()

	// write the counter value to stdout
	registry.Report(StdoutReporter)
}

func readMemoryUsageInMB() float64 {
	return 0
}

func ExampleGauge() {
	registry := NewRegistry("my-registry")
	registry.NewGauge("my-gauge", readMemoryUsageInMB)

	// write the gauge value to stdout
	registry.Report(StdoutReporter)
}

func ExampleTimer() {
	registry := NewRegistry("my-registry")
	timer := registry.NewTimer("my-timer", Milliseconds)

	stopwatch := timer.Start()
	// execute a task
	stopwatch.Record()

	// write the timer value to stout
	registry.Report(StdoutReporter)
}

func ExampleRegistry_report() {
	registry := NewRegistry("my-registry")
	// create and use some metrics

	registry.Report(StdoutReporter)
}

func ExampleReporting() {
	registry := NewRegistry("my-registry")

	// write the metrics of registry every five seconds to stdout
	reporting := StartReporting(5*time.Second, StdoutReporter)
	reporting.Attach(registry)

	// create and use some metric objects
}
