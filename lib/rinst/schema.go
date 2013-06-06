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

type SchemaBuffer chan Schema

// a schema describes a particular kind of measurement
type Schema struct {
	Name        string     // the name of an instrument
	Type        SchemaType // the kind of instrument
	Unit        string     // a description of the unit
	Description string     // a human readable description
}
