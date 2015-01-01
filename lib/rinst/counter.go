package rinst

import (
	"sync/atomic"
	"time"
)

// Counter is a monotonically increasing measurement, usually tracked through
// the first derivative.
type Counter struct {
	count      int64
	desc, unit string
}

// Incr atomically increments the given counter by one
func (c *Counter) Incr() *Counter {
	return c.Add(1)
}

// Add the given delta atomically to the counter
func (c *Counter) Add(del uint) *Counter {
	atomic.AddInt64(&c.count, int64(del))
	return c
}

// Value atomically gets the value of the counter.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.count)
}

// Measure gets a readable value for a counter.
func (c *Counter) Measure(name string, r MeasurementReceiver) {
	now := time.Now().Unix()
	r.WriteInt(name, now, c.Value())
}

// Schema writes the schema of the counter into the receiver.
func (c *Counter) Schema(name string, r SchemaReceiver) {
	r.WriteSchema(name, TypeCounter, c.unit, c.desc)
}
