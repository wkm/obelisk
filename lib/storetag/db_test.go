package storetag

import (
	_ "circuit/kit/debug/ctrlc"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestDB(t *testing.T) {
	c := NewConfig()
	c.DiskStore = filepath.Join(os.TempDir(), "obelisk-storetag")
	defer os.RemoveAll(c.DiskStore)

	db, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// insert some data
	db.Store.NewTag("a")
	db.Store.NewTag("b")
	db.Store.NewTag("c")
	db.Store.NewTag("d")
	db.Store.NewTag("e")
	db.Store.NewTag("g")
	db.Store.NewTag("h")

	db.Store.NewTag("a/i")
	db.Store.NewTag("a", "j")
	db.Store.NewTag("a/k")

	id1, _ := db.Store.NewTag("a/j/l")
	id2, _ := db.Store.NewTag("a/j/m")
	id3, _ := db.Store.NewTag("a/j/n")
	id4, _ := db.Store.NewTag("a/j/o")

	db.Flush()
	db.Shutdown()

	db2, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err := db2.Store.Children("a/j")
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}

	sort.Strings(res)
	str := strings.Join(res, ",")
	if str != "l,m,n,o" {
		t.Fatalf("expected l,m,n,o got %v", str)
	}

	r, _ := db.Store.Id("a/j/l")
	if r != id1 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db.Store.Id("a/j/m")
	if r != id2 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db.Store.Id("a/j/n")
	if r != id3 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	r, _ = db.Store.Id("a/j/o")
	if r != id4 {
		t.Errorf("expected id=%v got %v", id1, r)
	}

	db2.Shutdown()
}
