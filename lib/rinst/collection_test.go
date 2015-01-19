package rinst

import (
	"testing"
)

func TestLayout(t *testing.T) {
	t.Parallel()

	coll := NewCollection()
	c := coll.Counter("foo", "bar", "foos to the bar")
	c.Incr()
	c.Incr()

	if c.Value() != 2 {
		t.Errorf("layout counter bad value %d", c.Value())
	}

	b := make(MeasurementBuffer, 0)
	coll.Snapshot(&b)

	if len(b) != 1 {
		t.Errorf("layout snapshot bad length %d", len(b))
	}

	recv := b[0]
	if recv.Name != "foo" || recv.IntValue != 2 {
		t.Errorf("layout snapshot bad value %#v", recv)
	}
}
