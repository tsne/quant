package quant

type metric struct {
	name string
	unit string
}

// Name returns the name of the metric.
func (m *metric) Name() string {
	return m.name
}

// Unit returns a unit represenation of the metric or
// an empty string if no unit is specified.
func (m *metric) Unit() string {
	return m.unit
}
