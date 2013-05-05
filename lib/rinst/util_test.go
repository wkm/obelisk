package rinst

import (
	"fmt"
	"testing"
)

func TestFlushSchema(t *testing.T) {
	coll := NewCollection()

	for i := 0; i < 100; i++ {
		c := coll.Counter(fmt.Sprintf("%d", i), "unit", "desc")
		c.Incr()
	}

	sb := FlushSchema(coll, 3)
	if len(sb) != 100 {
		t.Errorf("expected buffer len=100 got %d", len(sb))
	}

	mb := FlushMeasurements(coll, 3)
	if len(mb) != 100 {
		t.Errorf("expected buffer len=100 got %d", len(mb))
	}
}
