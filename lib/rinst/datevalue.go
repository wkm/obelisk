package rinst

import (
	"fmt"
	"sync/atomic"
	"time"
)

// a datevalue stores a uint64 which represents a date
type DateValue struct {
	value   uint64
	changes uint32
	desc    string
}

// atomically set the value
func (v *DateValue) Set(value uint64) *DateValue {
	atomic.StoreUint64(&v.value, value)
	atomic.AddUint32(&v.changes, 1)
	return v
}

// atomically get the value of this value
func (v *DateValue) Get() uint64 {
	return atomic.LoadUint64(&v.value)
}

// atomically get the number of changes for this value
func (v *DateValue) NumSets() uint32 {
	return atomic.LoadUint32(&v.changes)
}

// get a readable value for a counter
func (v *DateValue) Measure(n string, b MeasurementBuffer) {
	now := uint64(time.Now().Unix())
	b <- Measurement{n, now, fmt.Sprintf("%d", v.Get())}
	b <- Measurement{n + ".sets", now, fmt.Sprintf("%d", v.NumSets())}
}

// the schema of this value
func (v *DateValue) Schema(name string, b SchemaBuffer) {
	b <- Schema{name, TypeDateValue, "", v.desc}
	b <- Schema{name + ".sets", TypeCounter, "set", "rate of changes to this value"}
}
