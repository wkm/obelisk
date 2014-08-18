package storetag

import (
	_ "circuit/kit/debug/ctrlc"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDB(t *testing.T) {
	c := Config{}
	c.DiskStore = filepath.Join(os.TempDir(), "github.com/wkm/obelisk-storetag")
	defer os.RemoveAll(c.DiskStore)

	db, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// insert some data
	db.NewTag("a")
	db.NewTag("b")
	db.NewTag("c")
	db.NewTag("d")
	db.NewTag("e")
	db.NewTag("g")
	db.NewTag("h")

	db.NewTag("a/i")
	db.NewTag("a", "j")
	db.NewTag("a/k")

	id1, _ := db.NewTag("a/j/l")
	id2, _ := db.NewTag("a/j/m")
	id3, _ := db.NewTag("a/j/n")
	id4, _ := db.NewTag("a/j/o")

	// Close and open the database to test persistence
	db.Shutdown()
	db2, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err := db2.Children("a/j")
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}

	str := strings.Join(res, ",")
	exp := "a/j,a/j/l,a/j/m,a/j/n,a/j/o"
	if str != exp {
		t.Fatalf("expected %#v got %#v", exp, str)
	}

	r, _ := db2.Id("a/j/l")
	if r != id1 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db2.Id("a/j/m")
	if r != id2 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db2.Id("a/j/n")
	if r != id3 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db2.Id("a/j/o")
	if r != id4 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	db2.Shutdown()
}
