package quant

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type reportingSettings struct {
	interval  time.Duration
	reporters []Reporter
}

// Reporting represents the periodic process of writing a set of metrics
// to specified locations.
type Reporting struct {
	settingsChan chan *reportingSettings
	wg           sync.WaitGroup
	mtx          sync.RWMutex
	registries   map[*Registry]struct{}
}

// StartReporting creates a new reporting with the specified interval and reporters.
// If the interval is not positive this function will panic.
func StartReporting(interval time.Duration, reporters ...Reporter) *Reporting {
	if interval <= 0 {
		panic(fmt.Errorf("invalid reporting interval: %s", interval.String()))
	}

	reporting := &Reporting{
		settingsChan: make(chan *reportingSettings),
		registries:   make(map[*Registry]struct{}),
	}

	reporting.wg.Add(1)
	go reporting.run(&reportingSettings{
		interval:  interval,
		reporters: reporters,
	})
	return reporting
}

// Attach adds a new registry to the reporting. If the registry
// was already attached to the reporting this function will panic.
func (r *Reporting) Attach(registry *Registry) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, exists := r.registries[registry]; exists {
		panic(fmt.Errorf("registry already exists in this reporting: %s", registry.Name()))
	}
	r.registries[registry] = struct{}{}
}

// Detach removes a registry from the reporting. If the registry
// was not attached to the reporting, this function is a no-op.
func (r *Reporting) Detach(registry *Registry) {
	r.mtx.Lock()
	delete(r.registries, registry)
	r.mtx.Unlock()
}

// Reset changes the interval and the reporters for the reporting.
// Calling this function is allowed on running reportings only. If the
// reporting was stopped this function will panic.
func (r *Reporting) Reset(interval time.Duration, reporters ...Reporter) {
	r.checkRunning()

	if interval <= 0 {
		panic(fmt.Errorf("invalid reporting interval: %s", interval.String()))
	}
	r.settingsChan <- &reportingSettings{
		interval:  interval,
		reporters: reporters,
	}
}

// Stop stops the reporting. Once the reporting is stopped
// no more metrics are written to the configured reporters.
// If there is a reporting running, Stop will wait until it
// is finished.
func (r *Reporting) Stop() {
	r.checkRunning()

	close(r.settingsChan)
	r.settingsChan = nil
	r.wg.Wait()
}

func (r *Reporting) allRegistries() []*Registry {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	registries := make([]*Registry, len(r.registries))
	idx := 0
	for registry, _ := range r.registries {
		registries[idx] = registry
		idx++
	}
	return registries
}

func (r *Reporting) checkRunning() {
	if r.settingsChan == nil {
		panic("reporting was stopped")
	}
}

func (r *Reporting) run(settings *reportingSettings) {
	defer r.wg.Done()

	ticker := time.NewTicker(settings.interval)
	defer ticker.Stop()

	for {
		select {
		case s, ok := <-r.settingsChan:
			if !ok {
				return
			}
			// restart the ticker
			ticker.Stop()
			ticker = time.NewTicker(s.interval)
			settings = s

		case <-ticker.C:
			for _, registry := range r.allRegistries() {
				err := registry.Report(settings.reporters...)
				if err != nil {
					log.Fatalf("error reporting registry %s: %s", registry.Name(), err)
				}
			}
		}
	}
}
