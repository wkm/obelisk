package rinst

import (
	"sync"
	"time"
)

// a container for an integer value which sometimes changes
type FloatValue struct {
	sync.Mutex
	value      float64
	desc, unit string
}

// atomically set the value
func (v *FloatValue) Set(value float64) *FloatValue {
	v.Lock()
	defer v.Unlock()

	v.value = value
	return v
}

// atomically get the value of this value
func (v *FloatValue) Get() float64 {
	v.Lock()
	defer v.Unlock()
	return v.value
}

// get a readable value for a counter
func (v *FloatValue) Measure(n string, r MeasurementReceiver) {
	v.Lock()
	defer v.Unlock()

	now := time.Now().Unix()
	r.WriteFloat(n, now, v.value)
}

// the schema of this value
func (v *FloatValue) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeFloatValue, v.unit, v.desc)
}
