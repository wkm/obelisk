package rinst

import (
	"testing"
)

// test that counters increment correctly
func TestCounter(t *testing.T) {
	c := &Counter{}
	c.Incr()
	c.Incr()

	if c.Value() != 2 {
		t.Error("counter bad value")
	}

	if len(c.Measure()) != 1 {
		t.Error("counter bad measure length")
	}
	if c.Measure()[0] != "2" {
		t.Error("counter bad measure")
	}
}
