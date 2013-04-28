/*
	implement a selection of instruments for measuring program and
	worker behavior and performance

	const Stats = make(Collection)
*/
package rinst

// an instrumentation collection
type Collection struct {
	instruments map[string]Instrument
}

// create a new collection
func NewCollection() *Collection {
	c := new(Collection)
	c.instruments = make(map[string]Instrument)
	return c
}

// add a new instrument to the layout
func (coll Collection) AddInstrument(name string, inst Instrument) Instrument {
	coll.instruments[name] = inst
	return coll.instruments[name]
}

// create a new counter with the given name
func (coll Collection) Counter(name string) *Counter {
	counter := new(Counter)
	coll.instruments[name] = counter
	return counter
}

// create a new integer value with the given name
func (coll Collection) IntValue(name string) *IntValue {
	value := new(IntValue)
	coll.instruments[name] = value
	return value
}

// create a new float value with the given name
func (coll Collection) FloatValue(name string) *FloatValue {
	value := new(FloatValue)
	coll.instruments[name] = value
	return value
}

// send the current values of all instruments in a layout to a buffer
func (coll Collection) Snapshot(b MeasurementBuffer) {
	coll.Measure("", b)
}

// measure all instruments in this collection
func (coll Collection) Measure(name string, buff MeasurementBuffer) {
	prefix := ""
	if name != "" {
		prefix = name + "."
	}
	for key, inst := range coll.instruments {
		inst.Measure(prefix+key, buff)
	}
}
