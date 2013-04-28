package rinst

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTextualReport(t *testing.T) {
	coll := NewCollection()
	c := coll.Counter("foo")
	c.Incr()
	c.Incr()

	var buff bytes.Buffer
	TextReport(&buff, coll)

	if buff.String() != fmt.Sprintf("%s: %s\n", "foo", "2") {
		t.Errorf("invalid text export: %s", buff.String())
	}
}
