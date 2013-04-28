package storetag

import (
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

	db.Store.NewChildTag("i", "a")
	db.Store.NewChildTag("j", "a")
	db.Store.NewChildTag("k", "a")

	db.Store.NewChildTag("l", "j")
	db.Store.NewChildTag("m", "j")
	db.Store.NewChildTag("n", "j")
	db.Store.NewChildTag("o", "j")

	db.Flush()
	db.Shutdown()

	db2, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err := db2.Store.Children("j")
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
