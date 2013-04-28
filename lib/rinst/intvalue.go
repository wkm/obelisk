package rinst

import (
	"fmt"
	"sync/atomic"
	"time"
)

// a container for an integer value which sometimes changes
type IntValue struct {
	value      int64
	changes    uint32
	desc, unit string
}

// atomically set the value
func (v *IntValue) Set(value int64) *IntValue {
	atomic.StoreInt64(&v.value, value)
	atomic.AddUint32(&v.changes, 1)
	return v
}

// atomically get the value of this value
func (v *IntValue) Get() int64 {
	return atomic.LoadInt64(&v.value)
}

// atomically get the number of changes for this value
func (v *IntValue) NumSets() uint32 {
	return atomic.LoadUint32(&v.changes)
}

// get a readable value for a counter
func (v *IntValue) Measure(n string, b MeasurementBuffer) {
	now := uint64(time.Now().Unix())
	b <- Measurement{n, now, fmt.Sprintf("%d", v.Get())}
	b <- Measurement{n + ".sets", now, fmt.Sprintf("%d", v.NumSets())}
}

// the schema of this value
func (v *IntValue) Schema(name string, b SchemaBuffer) {
	b <- Schema{name, TypeValue, v.unit, v.desc}
	b <- Schema{name + ".sets", TypeCounter, "set", "rate of changes to this value"}
}
