package rinst

import (
	"testing"
)

func TestLayout(t *testing.T) {
	coll := NewCollection()
	c := coll.Counter("foo")
	c.Incr()
	c.Incr()

	if c.Value() != 2 {
		t.Errorf("layout counter bad value %d", c.Value())
	}

	b := make(MeasurementBuffer, 10)
	coll.Snapshot(b)

	if len(b) != 1 {
		t.Errorf("layout snapshot bad length %d", len(b))
	}

	recv := <-b
	if recv.Name != "foo" || recv.Value != "2" {
		t.Error("layout snapshot bad value %#v", recv)
	}
}
