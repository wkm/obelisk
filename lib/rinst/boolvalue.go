package rinst

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// a boolean value represents a
type BoolValue struct {
	value      int32
	changes    uint32
	desc, unit string
	sync.Mutex
}

func (v *BoolValue) Set(b bool) *BoolValue {
	if b {
		atomic.StoreInt32(&v.value, 1)
	} else {
		atomic.StoreInt32(&v.value, 0)
	}
	return v
}

func (v *BoolValue) Get() bool {
	if atomic.LoadInt32(&v.value) == 0 {
		return false
	}
	return true
}

func (v *BoolValue) NumSets() uint32 {
	return atomic.LoadUint32(&v.changes)
}

func (v *BoolValue) Measure(n string, b MeasurementBuffer) {
	now := uint64(time.Now().Unix())
	b <- Measurement{n, now, fmt.Sprintf("%d", v.Get())}
	b <- Measurement{n + ".sets", now, fmt.Sprintf("%d", v.NumSets())}
}

func (v *BoolValue) Schema(name string, b SchemaBuffer) {
	b <- Schema{name, TypeBoolValue, v.unit, v.desc}
	b <- Schema{name + ".sets", TypeCounter, "set", "rate of changes to this value"}
}
