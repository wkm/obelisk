package rinst

// SchemaType represents the type of measurement.
type SchemaType uint

// Types of measurements.
const (
	TypeIntValue = iota
	TypeBoolValue
	TypeFloatValue
	TypeDateValue // An int representing a date value
	TypeCounter   // A monotonic counter
	TypeAllocation

	TypeFloatStream
)

// The SchemaReceiver interface specifies how an instrument can publish its schema.
type SchemaReceiver interface {
	WriteSchema(name string, ty SchemaType, unit, desc string)
}

// The MeasurementReceiver interface specifies how an instrument can publish its values.
type MeasurementReceiver interface {
	// WriteInt allows publishing integer values
	WriteInt(name string, time int64, value int64)

	// WriteFloat allows publishing floating point values.
	WriteFloat(name string, time int64, value float64)
}
