package rinst

type SchemaType uint

const (
	TypeIntValue = iota
	TypeBoolValue
	TypeFloatValue
	TypeDateValue // an int representing a date value
	TypeCounter   // a monotonic counter
	TypeAllocation

	TypeFloatStream
)

type SchemaReceiver interface {
	WriteSchema(name string, ty SchemaType, unit, desc string)
}

type MeasurementReceiver interface {
	WriteInt(name string, time uint64, value int64)
	WriteFloat(name string, time uint64, value float64)
}
