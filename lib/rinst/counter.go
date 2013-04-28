package rinst

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Counter struct {
	count      int64
	desc, unit string
}

// atomically increment the given counter by one
func (c *Counter) Incr() *Counter {
	return c.Add(1)
}

// atomically add the given delta to the counter
func (c *Counter) Add(del uint) *Counter {
	atomic.AddInt64(&c.count, int64(del))
	return c
}

// atomically get the value of the given counter
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.count)
}

// get a readable value for a counter
func (c *Counter) Measure(name string, b MeasurementBuffer) {
	now := uint64(time.Now().Unix())
	b <- Measurement{name, now, fmt.Sprintf("%d", c.Value())}
}

// the schema of this counter
func (c *Counter) Schema(name string, b SchemaBuffer) {
	b <- Schema{name, TypeCounter, c.unit, c.desc}
}
