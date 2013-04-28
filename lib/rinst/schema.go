package rinst

type SchemaType uint

const (
	TypeValue   = iota // a raw reading
	TypeCounter        // a monotonic counter
)

type SchemaBuffer chan Schema

// a schema describes a particular kind of measurement
type Schema struct {
	Name        string     // the name of an instrument
	Type        SchemaType // the kind of instrument
	Unit        string     // a description of the unit
	Description string     // a human readable description
}
