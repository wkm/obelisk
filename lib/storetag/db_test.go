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

	println("about to insert data")

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

	db.Store.NewTag("a/j/l")
	db.Store.NewTag("a/j/m")
	db.Store.NewTag("a/j/n")
	db.Store.NewTag("a/j/o")

	println("about to enter flush")
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

	db2.Shutdown()
}
