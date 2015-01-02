package rinst

import (
	"sync/atomic"
	"time"
)

// IntValue stores a int64 as an instrument.
type IntValue struct {
	value      int64
	desc, unit string
}

// Set the value of this instrument atomically.
func (v *IntValue) Set(value int64) *IntValue {
	atomic.StoreInt64(&v.value, value)
	return v
}

// Get the value of this instrument atomically.
func (v *IntValue) Get() int64 {
	return atomic.LoadInt64(&v.value)
}

// Measure the value of instrument into the receiver.
func (v *IntValue) Measure(n string, r MeasurementReceiver) {
	now := time.Now().Unix()
	r.WriteInt(n, now, v.Get())
}

// Schema writes the schema of this value into the receiver.
func (v *IntValue) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeIntValue, v.unit, v.desc)
}
