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

	b := make(MeasurementBuffer, 0)
	c.Measure("fig", &b)

	if len(b) != 1 {
		t.Error("counter bad measure length")
	}

	recv := b[0]
	if recv.Name != "fig" || recv.IntValue != 2 {
		t.Errorf("counter bad measure: %#v", recv)
	}
}
