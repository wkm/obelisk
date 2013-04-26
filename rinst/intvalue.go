package rinst

import (
	"fmt"
	"sync/atomic"
)

// a container for an integer value which sometimes changes
type IntValue struct {
	value   int64
	changes uint32
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
func (v *IntValue) NumChanges() uint32 {
	return atomic.LoadUint32(&v.changes)
}

// get a readable value for a counter
func (c *IntValue) Measure() []string {
	return []string{
		fmt.Sprintf("%d", c.Get()),
		fmt.Sprintf("%d", c.NumChanges()),
	}
}
