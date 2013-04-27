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

	b := make(MeasurementBuffer, 10)
	c.Measure("fig", b)

	if len(b) != 1 {
		t.Error("counter bad measure length")
	}

	expec := Measurement{"fig", "2"}
	if expec != <-b {
		t.Error("counter bad measure")
	}
}
