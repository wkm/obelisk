package rinst

type MeasurementBuffer chan Measurement

// represents a single value measured by an instrument
type Measurement struct {
	Name  string
	Value string
}
