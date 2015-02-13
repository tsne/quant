package quant

import (
	"math"
)

type snapshot struct {
	name string
	unit string
}

// Name returns the name of the metric this snapshot
// belongs to.
func (s *snapshot) Name() string {
	return s.name
}

// Unit returns the unit representation of the mtric this
// snapshot belongs to. If the metric has no unit an empty
// string will be returned.
func (s *snapshot) Unit() string {
	return s.unit
}

type reservoirSnapshot struct {
	snapshot
	count int
	min   float64
	max   float64
	sum   float64
	sumSq float64
}

func newReservoirSnaphot(name, unit string) *reservoirSnapshot {
	return &reservoirSnapshot{
		snapshot: snapshot{name, unit},
		count:    0,
		min:      0,
		max:      0,
		sum:      0,
		sumSq:    0,
	}
}

// Count returns the number of measurements this snapshot contains.
func (s *reservoirSnapshot) Count() int {
	return s.count
}

// Minimum returns the smallest value this snapshot contains.
func (s *reservoirSnapshot) Minimum() float64 {
	return s.min
}

// Maximum returns the biggest value this snapshot contains.
func (s *reservoirSnapshot) Maximum() float64 {
	return s.max
}

// Average returns the mean of all values this snapshot contains.
func (s *reservoirSnapshot) Average() float64 {
	return s.sum / float64(s.count)
}

// Variance returns the variance of all values this snapshot contains.
func (s *reservoirSnapshot) Variance() float64 {
	return (s.sumSq - s.sum*s.sum/float64(s.count)) / float64(s.count)
}

// StdDeviation returns the standard deviation of all value this snapshot
// contains. The result is equivalent to the square root of Variance.
func (s *reservoirSnapshot) StdDeviation() float64 {
	return math.Sqrt(s.Variance())
}

func (s *reservoirSnapshot) add(value float64) {
	if s.count == 0 {
		s.min = value
		s.max = value
	} else {
		if value < s.min {
			s.min = value
		}
		if value > s.max {
			s.max = value
		}
	}

	s.count++
	s.sum += value
	s.sumSq += value * value
}
