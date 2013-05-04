package rinst

// a container for values which are measured upon request
type GaugeValue struct {
	MeasureFn func(name string, b MeasurementBuffer)
	SchemaFn  func(name string, b SchemaBuffer)
}

func (g *GaugeValue) Measure(name string, b MeasurementBuffer) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	g.MeasureFn(prefix, b)
}

func (g *GaugeValue) Schema(name string, b SchemaBuffer) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	g.SchemaFn(prefix, b)
}
