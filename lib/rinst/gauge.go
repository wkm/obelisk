package rinst

// GaugeValue is a container for values which are measured during their read. The schema of a gauge is queried on initialization, and should not change over time.
type GaugeValue struct {
	MeasureFn func(name string, r MeasurementReceiver)
	SchemaFn  func(name string, r SchemaReceiver)
}

// Measure runs the measurement function and stores the values into the receiver.
func (g *GaugeValue) Measure(name string, r MeasurementReceiver) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	g.MeasureFn(prefix, r)
}

// Schema writes the schema of this gauge into the receiever.
func (g *GaugeValue) Schema(name string, r SchemaReceiver) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	g.SchemaFn(prefix, r)
}
