package storetime

import (
	"bytes"
	"fmt"
	"testing"
)

// functional test of the storetime's features
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

	db2 := NewStore()
	db2.Load(&b)

	var b2 bytes.Buffer
	db2.Dump(&b2)
	expec, _ := db.FlatQuery(123, 10, 12)
	actual, _ := db2.FlatQuery(123, 10, 12)
	if fmt.Sprintf("%v", expec) != fmt.Sprintf("%v", actual) {
		t.Errorf("queries are different")
	}
}
