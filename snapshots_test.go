package quant

import (
	"testing"
)

func TestReservoirSnapshotMinMaxAvg(t *testing.T) {
	s := newReservoirSnaphot("reservoir", "")
	s.add(2)
	s.add(1)
	s.add(3)

	if s.Count() != 3 {
		t.Errorf("wrong reservoir snapshot count: %d (3 expected)", s.Count())
	}
	if s.Minimum() != 1 {
		t.Errorf("wrong reservoir snapshot minimum: %f (1 expected)", s.Minimum())
	}
	if s.Maximum() != 3 {
		t.Errorf("wrong reservoir snapshot maximum: %f (3 expected)", s.Maximum())
	}
	if s.Average() != 2 {
		t.Errorf("wrong reservoir snapshot average: %f (2 expected)", s.Average())
	}
}

func TestReservoirSnapshotVariance(t *testing.T) {
	s := newReservoirSnaphot("reservoir", "")
	s.add(9)
	s.add(2)
	s.add(5)
	s.add(4)
	s.add(12)
	s.add(7)
	s.add(8)
	s.add(11)
	s.add(9)
	s.add(3)
	s.add(7)
	s.add(4)
	s.add(12)
	s.add(5)
	s.add(4)
	s.add(10)
	s.add(9)
	s.add(6)
	s.add(9)
	s.add(4)

	const expected = 178.0 / 20.0
	if s.Variance() != expected {
		t.Errorf("wrong reservoir snapshot variance: %f (%f expected)", s.Variance(), expected)
	}
}
