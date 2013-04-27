/*
	implement a selection of instruments for measuring program and
	worker behavior and performance

	const Stats = make(Collection)
*/
package rinst

// an instrumentation collection
type Collection map[string]Instrument

// add a new instrument to the layout
func (coll Collection) AddInstrument(name string, inst Instrument) Instrument {
	coll[name] = inst
	return coll[name]
}

// create a new counter with the given name
func (coll Collection) Counter(name string) *Counter {
	counter := new(Counter)
	coll[name] = counter
	return counter
}

// create a new integer value with the given name
func (coll Collection) IntValue(name string) *IntValue {
	value := new(IntValue)
	coll[name] = value
	return value
}

// create a new float value with the given name
func (coll Collection) FloatValue(name string) *FloatValue {
	value := new(FloatValue)
	coll[name] = value
	return value
}

// send the current values of all instruments in a layout to a buffer
func (coll Collection) Snapshot(b MeasurementBuffer) {
	for name, i := range coll {
		i.Measure(name, b)
	}
}
