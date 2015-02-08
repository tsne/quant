package quant

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	c := newCounter("my-counter", "")

	if c.Name() != "my-counter" {
		t.Errorf("wrong counter name: %s", c.Name())
	}
	if c.Value() != 0 {
		t.Errorf("wrong counter value: %d (0 expected)", c.Value())
	}

	c.Increment()
	c.Increment()
	c.Decrement()
	if c.Value() != 1 {
		t.Errorf("wrong counter value: %d (1 expected)", c.Value())
	}

	c.Reset()
	if c.Value() != 0 {
		t.Errorf("wrong counter value: %d (0 expected)", c.Value())
	}
}

func TestConcurrentCounter(t *testing.T) {
	const loops = 1000000
	var wg sync.WaitGroup
	c := newCounter("my-counter", "")

	wg.Add(2)
	go func() {
		for i := 0; i < loops; i++ {
			c.Increment()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < loops; i++ {
			c.Decrement()
		}
		wg.Done()
	}()

	wg.Wait()
	if c.Value() != 0 {
		t.Errorf("wrong counter value: %d (0 expected)", c.Value())
	}
}
