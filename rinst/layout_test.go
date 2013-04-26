package rinst

import (
	"testing"
)

func TestLayout(t *testing.T) {
	l := NewLayout()
	c := l.Counter("foo")
	c.Incr()
	c.Incr()

	if c.Value() != 2 {
		t.Errorf("layout counter bad value %d", c.Value())
	}

	snap := l.Snapshot()
	t.Logf("snaphost: %d", snap)
	if len(snap) != 1 {
		t.Errorf("layout snapshot bad length %d", len(snap))
	}

	expected := InstrumentMeasurement{"foo", "2"}
	if snap[0] != expected {
		t.Error("layout snapshot bad value")
	}
}
