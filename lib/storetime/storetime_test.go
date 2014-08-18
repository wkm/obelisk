package storetime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKey(t *testing.T) {
	if string(createKey(123, 345)) != "123•345" {
		t.Errorf("invalid key created")
	}

	if getTime([]byte("123•345")) != 345 {
		t.Errorf("invalid time extracted")
	}
}

func TestDB(t *testing.T) {
	c := Config{}
	c.DiskStore = filepath.Join(os.TempDir(), "obelisk-storetime")
	defer os.RemoveAll(c.DiskStore)

	db, err := NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// insert a bunch of data
	for i := uint64(0); i < 100; i++ {
		for j := uint64(0); j < 50; j++ {
			db.Insert(100+i, 10+j, float64(j))
		}
	}

	// validate data
	for i := uint64(0); i < 100; i++ {
		points, err := db.FlatQuery(100+i, 10, 10+50)
		if err != nil {
			t.Fatalf("unexpected error %s", err.Error())
		}

		if len(points) != 50 {
			t.Fatalf("expected 50 points, had %d", len(points))
		}

		for j := 0; j < 50; j++ {
			expec := float64(j)
			if points[j] != expec {
				t.Errorf("invalid point value %v, expected %v", points[j], expec)
			}
		}
	}

	// make sure we can't create another db
	db2, err := NewDB(c)
	if err == nil {
		t.Error("second DB created")
	}

	// close the original database
	db.Shutdown()

	// now reopen it
	db2, err = NewDB(c)
	if err != nil {
		t.Fatal(err.Error())
	}

	// validate data
	for i := uint64(0); i < 100; i++ {
		points, err := db2.FlatQuery(100+i, 10, 10+50)
		if err != nil {
			t.Fatalf("unexpected error %s", err.Error())
		}

		if len(points) != 50 {
			t.Fatalf("expected 50 points, had %d", len(points))
		}

		for j := 0; j < 50; j++ {
			expec := float64(j)
			if points[j] != expec {
				t.Fatalf("invalid point value %v, expected %v", points[j], expec)
			}
		}
	}

	db2.Shutdown()
}
