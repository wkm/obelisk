package rinst

import (
	"fmt"
	"sync"
	"time"
)

// a container for an integer value which sometimes changes
type FloatValue struct {
	sync.Mutex
	value      float64
	changes    uint32
	desc, unit string
}

// atomically set the value
func (v *FloatValue) Set(value float64) *FloatValue {
	v.Lock()
	defer v.Unlock()

	v.value = value
	v.changes++
	return v
}

// atomically get the value of this value
func (v *FloatValue) Get() float64 {
	v.Lock()
	defer v.Unlock()
	return v.value
}

// atomically get the number of changes for this value
func (v *FloatValue) NumSets() uint32 {
	v.Lock()
	defer v.Unlock()
	return v.changes
}

// get a readable value for a counter
func (v *FloatValue) Measure(n string, b MeasurementBuffer) {
	v.Lock()
	defer v.Unlock()

	now := uint64(time.Now().Unix())
	b <- Measurement{n, now, fmt.Sprintf("%f", v.value)}
	b <- Measurement{n + ".sets", now, fmt.Sprintf("%d", v.changes)}
}

// the schema of this value
func (v *FloatValue) Schema(name string, b SchemaBuffer) {
	b <- Schema{name, TypeValue, v.unit, v.desc}
	b <- Schema{name + ".sets", TypeCounter, "set", "rate of changes to this value"}
}
