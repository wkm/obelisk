package storekv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
	t.Parallel()

	c := Config{}
	c.DiskStore = filepath.Join(os.TempDir(), "obelisk-storekv")
	defer os.RemoveAll(c.DiskStore)

	db, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// shove some data in
	db.Set("a", []byte("1"))
	db.Set("a", []byte("1.1"))
	db.Set("b", []byte("3"))

	// test a value
	bb, err := db.Get("a")
	if string(bb) != "1.1" {
		t.Errorf("Expected %s, got %s", "1", bb)
	}

	// Test for locking
	db2, err := NewDB(c)
	if err == nil {
		t.Error("second DB created")
	}

	db.Shutdown()

	db2, err = NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual, err := db2.Get("a")
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != "1.1" {
		t.Errorf("expected 1.1, got %v", actual)
	}

	actual, err = db2.Get("b")
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != "3" {
		t.Errorf("expected 3, got %v", actual)
	}

	db2.Shutdown()
}
