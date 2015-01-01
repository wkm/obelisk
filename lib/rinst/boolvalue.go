package rinst

import (
	"sync/atomic"
	"time"
)

// BoolValue stores an integer which can take on two values as an instrument.
type BoolValue struct {
	value      int32
	desc, unit string
}

// Set the value of this instrument atomically.
func (v *BoolValue) Set(b bool) *BoolValue {
	if b {
		atomic.StoreInt32(&v.value, 1)
	} else {
		atomic.StoreInt32(&v.value, 0)
	}
	return v
}

// Get the value of this instrument atomically.
func (v *BoolValue) Get() bool {
	if atomic.LoadInt32(&v.value) == 0 {
		return false
	}
	return true
}

// Measure the value of instrument into the receiver.
func (v *BoolValue) Measure(n string, b MeasurementReceiver) {
	i := atomic.LoadInt32(&v.value)
	now := time.Now().Unix()
	b.WriteInt(n, now, int64(i))
}

// Schema writes the schema of this value into the receiver.
func (v *BoolValue) Schema(name string, b SchemaReceiver) {
	b.WriteSchema(name, TypeBoolValue, v.unit, v.desc)
}
