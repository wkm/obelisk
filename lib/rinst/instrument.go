package rinst

// Instrument represents a metric which has a schema and can be measured.
type Instrument interface {
	// Measure writes the instrument's measurement buffer
	Measure(name string, r MeasurementReceiver)

	// Schema gets the instrument's schema
	Schema(name string, r SchemaReceiver)
}
