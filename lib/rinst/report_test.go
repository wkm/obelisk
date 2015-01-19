package rinst

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTextualReport(t *testing.T) {
	t.Parallel()

	coll := NewCollection()
	c := coll.Counter("foo", "bar", "foos of bar")
	c.Incr()
	c.Incr()

	var buff bytes.Buffer
	TextReport(&buff, coll)

	exp := fmt.Sprintf("%s: %d %f\n", "foo", 2, 0.0)
	if buff.String() != exp {
		t.Errorf("invalid text export:\ngot: %#v\nexp: %#v", buff.String(), exp)
	}
}
