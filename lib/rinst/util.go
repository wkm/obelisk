package rinst

// InstrumentSchema contains metadata around a single instance of an instrument.
type InstrumentSchema struct {
	Name        string
	Type        SchemaType
	Unit        string
	Description string
}

// SchemaBuffer is a simple array-backed receiver for instrument schemas.
type SchemaBuffer []InstrumentSchema

// WriteSchema stores metadata into a SchemaBuffer.
func (sb *SchemaBuffer) WriteSchema(name string, ty SchemaType, unit, desc string) {
	*sb = append(*sb, InstrumentSchema{name, ty, unit, desc})
}

// InstrumentMeasurement contains data from a single measurement of an instrument.
type InstrumentMeasurement struct {
	Name       string
	Time       int64
	IntValue   int64
	FloatValue float64
}

// MeasurementBuffer is a simple array-backed receiver for instrument measurements.
type MeasurementBuffer []InstrumentMeasurement

// WriteInt stores an int in the measuremement buffer.
func (mb *MeasurementBuffer) WriteInt(name string, time int64, value int64) {
	*mb = append(*mb, InstrumentMeasurement{name, time, value, 0})
}

// WriteFloat stores a float in the measurement buffer.
func (mb *MeasurementBuffer) WriteFloat(name string, time int64, value float64) {
	*mb = append(*mb, InstrumentMeasurement{name, time, 0, value})
}
