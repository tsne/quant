package quant

import (
	"testing"
)

func TestGauge(t *testing.T) {
	g := newGauge("my-gauge", "", func() float64 { return 7.0 })

	if g.Name() != "my-gauge" {
		t.Errorf("wrong gauge name: %s", g.Name())
	}
	if g.Value() != 7.0 {
		t.Errorf("wrong gauge value: %f", g.Value())
	}
}
