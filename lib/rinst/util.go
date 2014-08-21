package rinst

type InstrumentSchema struct {
	Name        string
	Type        SchemaType
	Unit        string
	Description string
}

type SchemaBuffer []InstrumentSchema

func (sb *SchemaBuffer) WriteSchema(name string, ty SchemaType, unit, desc string) {
	*sb = append(*sb, InstrumentSchema{name, ty, unit, desc})
}

type InstrumentMeasurement struct {
	Name       string
	Time       int64
	IntValue   int64
	FloatValue float64
}

type MeasurementBuffer []InstrumentMeasurement

func (mb *MeasurementBuffer) WriteInt(name string, time int64, value int64) {
	*mb = append(*mb, InstrumentMeasurement{name, time, value, 0})
}

func (mb *MeasurementBuffer) WriteFloat(name string, time int64, value float64) {
	*mb = append(*mb, InstrumentMeasurement{name, time, 0, value})
}
