package rinst

// an instrument can be measured
type Instrument interface {
	// write the instrument's measurement buffer
	Measure(name string, buff MeasurementBuffer)

	// get the instrument's schema
	Schema(name string, buff SchemaBuffer)
}
