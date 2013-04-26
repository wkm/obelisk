/*
	implement a selection of instruments for measuring program and
	worker behavior and performance

	const Inst = Layout&
*/
package rinst

// an instrumentation layout
type Layout struct {
	Instruments map[string]Instrument
}

func NewLayout() *Layout {
	l := new(Layout)
	l.Instruments = make(map[string]Instrument)
	return l
}

// add a new instrument to the layout
func (l *Layout) AddInstrument(name string, inst Instrument) Instrument {
	l.Instruments[name] = inst
	return l.Instruments[name]
}

// create a new namespace for instruments
// func (l *Layout) Namespace(name string) *Layout {
// 	namespace := Namespace{name}
// 	l.Layout[name] = namespace
// }

// create a new counter with the given name
func (l *Layout) Counter(name string) *Counter {
	counter := new(Counter)
	l.AddInstrument(name, counter)
	return counter
}

// send the current values of all instruments in a layout to a buffer
func (l *Layout) Snapshot(b MeasurementBuffer) {
	for name, i := range l.Instruments {
		i.Measure(name, b)
	}
}
