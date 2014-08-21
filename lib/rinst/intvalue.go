package rinst

import (
	"sync/atomic"
	"time"
)

// a container for an integer value which sometimes changes
type IntValue struct {
	value      int64
	desc, unit string
}

// atomically set the value
func (v *IntValue) Set(value int64) *IntValue {
	atomic.StoreInt64(&v.value, value)
	return v
}

// atomically get the value of this value
func (v *IntValue) Get() int64 {
	return atomic.LoadInt64(&v.value)
}

// get a readable value for a counter
func (v *IntValue) Measure(n string, r MeasurementReceiver) {
	now := time.Now().Unix()
	r.WriteInt(n, now, v.Get())
}

// the schema of this value
func (v *IntValue) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeIntValue, v.unit, v.desc)
}
