package rinst

import (
	"sync"
	"time"
)

// FloatValue stores a float64 as an instrument.
type FloatValue struct {
	sync.Mutex
	value      float64
	desc, unit string
}

// Set the value of this instrument atomically.
func (v *FloatValue) Set(value float64) *FloatValue {
	v.Lock()
	defer v.Unlock()

	v.value = value
	return v
}

// Get the value of this instrument atomically.
func (v *FloatValue) Get() float64 {
	v.Lock()
	defer v.Unlock()
	return v.value
}

// Measure the value of instrument into the receiver.
func (v *FloatValue) Measure(n string, r MeasurementReceiver) {
	v.Lock()
	defer v.Unlock()

	now := time.Now().Unix()
	r.WriteFloat(n, now, v.value)
}

// Schema writes the schema of this value into the receiver.
func (v *FloatValue) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeFloatValue, v.unit, v.desc)
}
