package storetag

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDB(t *testing.T) {
	c := Config{}
	c.DiskStore = filepath.Join(os.TempDir(), "obelisk-storetag")
	defer os.RemoveAll(c.DiskStore)

	db, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// insert some data
	db.Tag(0, "a")
	db.Tag(1, "b")
	db.Tag(2, "c")
	db.Tag(3, "d")
	db.Tag(4, "e")
	db.Tag(5, "g")
	db.Tag(6, "h")

	db.Tag(7, "a/i")
	db.Tag(8, "a", "j")
	db.Tag(9, "a/k")

	id1, _ := db.Tag(10, "a/j/l")
	id2, _ := db.Tag(11, "a/j/m")
	id3, _ := db.Tag(12, "a/j/n")
	id4, _ := db.Tag(13, "a/j/o")

	// Close and open the database to test persistence
	db.Shutdown()
	db2, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	id, err := db2.Tag(5, "a")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	if id != 0 {
		t.Errorf("expected ID to be %d as originally set not %d", 0, id)
	}

	res, err := db2.Children("a/j")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	str := strings.Join(res, ",")
	exp := "a/j,a/j/l,a/j/m,a/j/n,a/j/o"
	if str != exp {
		t.Errorf("expected %#v got %#v", exp, str)
	}

	r, _ := db2.ID("a/j/l")
	if r != id1 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db2.ID("a/j/m")
	if r != id2 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db2.ID("a/j/n")
	if r != id3 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db2.ID("a/j/o")
	if r != id4 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	db2.Shutdown()
}
