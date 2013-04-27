package storekv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
	c := NewConfig()
	c.DiskStore = filepath.Join(os.TempDir(), "obelisk-storekv")
	defer os.RemoveAll(c.DiskStore)

	db, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// shove some data in
	db.Store.Set("a", []byte("1"))
	db.Store.Set("a", []byte("1.1"))
	db.Store.Set("b", []byte("3"))

	db.Flush()
	db2, err := NewDB(c)
	if err == nil {
		t.Error("second DB created")
	}

	db.Shutdown()

	db2, err = NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual, err := db.Store.Get("a")
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != "1.1" {
		t.Errorf("expected 1.1, got %v", actual)
	}

	actual, err = db.Store.Get("b")
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != "3" {
		t.Errorf("expected 3, got %v", actual)
	}

	db2.Shutdown()
}
