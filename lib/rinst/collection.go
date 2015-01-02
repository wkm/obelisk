// Package rinst implement a selection of instruments for measuring  program and worker behavior and performance.
//	    const Stats = make(Collection)
package rinst

// Collection is a set of named instruments.
type Collection struct {
	instruments map[string]Instrument
}

// NewCollection creates a new set of named instruments.
func NewCollection() *Collection {
	c := new(Collection)
	c.instruments = make(map[string]Instrument)
	return c
}

// AddInstrument add a new instrument to the layout.
func (coll Collection) AddInstrument(name string, inst Instrument) Instrument {
	coll.instruments[name] = inst
	return coll.instruments[name]
}

// Counter creates a new counter with the given name.
func (coll Collection) Counter(name, unit, desc string) *Counter {
	counter := new(Counter)
	counter.unit = unit
	counter.desc = desc

	coll.instruments[name] = counter
	return counter
}

// IntValue creates a new integer value with the given name.
func (coll Collection) IntValue(name, unit, desc string) *IntValue {
	value := new(IntValue)
	value.unit = unit
	value.desc = desc

	coll.instruments[name] = value
	return value
}

// FloatValue creates a new float value with the given name.
func (coll Collection) FloatValue(name, unit, desc string) *FloatValue {
	value := new(FloatValue)
	value.unit = unit
	value.desc = desc

	coll.instruments[name] = value
	return value
}

// Snapshot sends the current values of all instruments into the receiever.
func (coll Collection) Snapshot(r MeasurementReceiver) {
	coll.Measure("", r)
}

// Measure all instruments in this collection.
func (coll Collection) Measure(name string, r MeasurementReceiver) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	for key, inst := range coll.instruments {
		inst.Measure(prefix+key, r)
	}
}

// Schema sends the schema of all instruments into the receiever.
func (coll Collection) Schema(name string, r SchemaReceiver) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	for key, inst := range coll.instruments {
		inst.Schema(prefix+key, r)
	}
}
