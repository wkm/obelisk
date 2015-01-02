package rinst

import (
	"sync/atomic"
	"time"
)

// DateValue stores a int64 which represents a date.
type DateValue struct {
	value int64
	desc  string
}

// Set the value of this instrument atomically.
func (v *DateValue) Set(value int64) *DateValue {
	atomic.StoreInt64(&v.value, value)
	return v
}

// Get the value of this instrument atomically.
func (v *DateValue) Get() int64 {
	return atomic.LoadInt64(&v.value)
}

// Measure the value of instrument into the receiver.
func (v *DateValue) Measure(n string, r MeasurementReceiver) {
	now := time.Now().Unix()
	r.WriteInt(n, now, v.Get())
}

// Schema writes the schema of this value into the receiver.
func (v *DateValue) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeDateValue, "", v.desc)
}
