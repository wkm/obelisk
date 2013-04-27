package timestore

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	db := NewStore()
	db.Insert(123, 10, 1.1)
	db.Insert(123, 11, 1.2)
	db.Insert(123, 12, 1.3)

	values, err := db.FlatQuery(123, 10, 12)
	if err != nil {
		t.Errorf("didn't expect error %s", err.Error())
	}
	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}

	var b bytes.Buffer
	db.Dump(&b)

	if b.String() != "123:10=1.100000,11=1.200000,12=1.300000,\n" {
		t.Errorf("invalid dump %s", b.String())
	}
}
