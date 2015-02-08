package quant

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	tm := newTimer("my-timer", Milliseconds)

	if tm.Name() != "my-timer" {
		t.Errorf("wrong timer name: %s", tm.Name())
	}
}

func TestStopwatch(t *testing.T) {
	tm := newTimer("my-timer", Milliseconds)
	sw := tm.Start()

	time.Sleep(10 * time.Millisecond)
	d := sw.Elapsed()
	if d < 10*time.Millisecond {
		t.Errorf("wrong time: %s (at least 10ms expected)", d)
	}
}
