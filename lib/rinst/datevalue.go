package rinst

import (
	"sync/atomic"
	"time"
)

// DateValue stores a int64 which represents a date
type DateValue struct {
	value int64
	desc  string
}

// atomically set the value
func (v *DateValue) Set(value int64) *DateValue {
	atomic.StoreInt64(&v.value, value)
	return v
}

// atomically get the value of this value
func (v *DateValue) Get() int64 {
	return atomic.LoadInt64(&v.value)
}

// get a readable value for a counter
func (v *DateValue) Measure(n string, r MeasurementReceiver) {
	now := time.Now().Unix()
	r.WriteInt(n, now, v.Get())
}

// the schema of this value
func (v *DateValue) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeDateValue, "", v.desc)
}
