[![GoDoc](https://godoc.org/github.com/thonit/quant?status.png)](https://godoc.org/github.com/thonit/quant)

# quant

quant is a simple metrics library for Go. It provides the following metric types to measure
application statistics:
* Counter
* Gauge
* Timer

Use `go get` to install or update the package:
```
go get -u github.com/thonit/quant
```

## Getting Started
The first step to use application metrics is to create a registry which acts as a
collection of metrics. With this registry all supported metric types can be created.
Each metric has its own unique name within the registry to identify the metric.
A registry provides the function `Report` to write a snapshot of each registered
metric to the specified reporters. A `Reporter` writes the snapshot to the specified
location in the specified format. The quant package provides the following reporters:
* `NullReporter`: does not write any snapshot
* `StdoutReporter`: writes the snapshots to the standard output

For a better metrics tracking snapshots of the metrics could be constantly written
to a specific location (e.g. a database). This can be achieved in two ways: Periodically
call the `Report` function of the registry, or starting a reporting and attach the
registry to it.

A complete example calling `Registry.Report`:
```go
package main

import (
	"runtime"
	"time"

	"github.com/thonit/quant"
)

func main() {
	registry := quant.NewRegistry("my-registry")

	go func() {
		for range time.Tick(time.Second) {
			registry.Report(quant.StdoutReporter)
		}
	}()

	// use the registry
	counter := registry.NewCounter("my-counter")
	timer := registry.NewTimer("my-timer", quant.Milliseconds)
	registry.NewGaugeWithUnit("my-gauge", "MB", readMemoryUsageInMB)

	for i := 0; i < 1000; i++ {
		counter.Increment()
		stopwatch := timer.Start()
		time.Sleep(10 * time.Millisecond) // simulate payload
		stopwatch.Record()
	}
}

func readMemoryUsageInMB() float64 {
	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)
	return float64(memstats.Alloc) / (1024.0 * 1024.0)
}
```

A complete example using `Reporting`:
```go
package main

import (
	"runtime"
	"time"

	"github.com/thonit/quant"
)

func main() {
	registry := quant.NewRegistry("my-registry")

	reporting := quant.StartReporting(time.Second, quant.StdoutReporter)
	defer reporting.Stop()
	reporting.Attach(registry)

	// use the registry
	counter := registry.NewCounter("my-counter")
	timer := registry.NewTimer("my-timer", quant.Milliseconds)
	registry.NewGaugeWithUnit("my-gauge", "MB", readMemoryUsageInMB)

	for i := 0; i < 1000; i++ {
		counter.Increment()
		stopwatch := timer.Start()
		time.Sleep(10 * time.Millisecond) // simulate payload
		stopwatch.Record()
	}
}

func readMemoryUsageInMB() float64 {
	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)
	return float64(memstats.Alloc) / (1024.0 * 1024.0)
}
```

## Supported Metrics
### Counters
A counter reports a single integral value. As the name says, it counts the occurence of
specific events and provides a functions to increment, decrement and reset the counter value.
These operations are thread-safe.

### Gauges
A gauge reports a single floating point value. The function providing this value is called
gauge reader is specified by the application. It is wrapped in a thread-safe context, i.e.
no extra synchronization is needed. Examples for using a gauge: reporting the memory
consumption or a buffer's size.

### Timers
A timer reports a series of measured time durations. When starting a timer a stopwatch
is created which immediately starts the measurement. Each stopwatch can report its measured
duration to the underlying timer. So a series of durations is created which could be
reported by the registry the timer belongs to. A stopwatch is not thread-safe and therefore
should not be used concurrently. A timer on the other hand is thread-safe.
