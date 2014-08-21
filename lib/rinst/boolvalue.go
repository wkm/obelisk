package rinst

import (
	"sync"
	"sync/atomic"
	"time"
)

// a boolean value represents a
type BoolValue struct {
	value      int32
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

func (v *BoolValue) Measure(n string, b MeasurementReceiver) {
	i := atomic.LoadInt32(&v.value)
	now := time.Now().Unix()
	b.WriteInt(n, now, int64(i))
}

func (v *BoolValue) Schema(name string, b SchemaReceiver) {
	b.WriteSchema(name, TypeBoolValue, v.unit, v.desc)
}
