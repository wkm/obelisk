package rinst

import (
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

func (v *BoolValue) Measure(n string, b MeasurementReceiver) {
	i := atomic.LoadInt32(&v.value)
	now := uint64(time.Now().Unix())
	b.WriteInt(n, now, i)
	b.WriteInt(n+".sets", now, v.NumSets())
}

func (v *BoolValue) Schema(name string, b SchemaReceiver) {
	b.WriteSchema(name, TypeBoolValue, unit, desc)
	b.WriteSchema(name+".sets", TypeCounter, "set", "rate of change on this value")
}
