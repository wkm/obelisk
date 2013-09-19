package rinst

type MeasurementBuffer chan Measurement

// A single value measured by an instrument
type Measurement struct {
	Name  string
	Time  uint64
	Value string
}
