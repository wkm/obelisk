package rinst

// a container for values which are measured upon request
type GaugeValue struct {
	MeasureFn func(name string, r MeasurementReceiver)
	SchemaFn  func(name string, r SchemaReceiver)
}

func (g *GaugeValue) Measure(name string, r MeasurementReceiver) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	g.MeasureFn(prefix, r)
}

func (g *GaugeValue) Schema(name string, r SchemaReceiver) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	g.SchemaFn(prefix, r)
}
