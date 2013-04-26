package rinst

import (
	"fmt"
	"sync/atomic"
)

type Counter struct {
	count int64
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
func (c *Counter) Measure() []string {
	return []string{
		fmt.Sprintf("%d", c.Value()),
	}
}
