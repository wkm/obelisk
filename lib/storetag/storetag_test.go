package storetag

import (
	"testing"
)

func TestTagStore(t *testing.T) {
	s := NewStore()

	if s.MaxId() != 1 {
		t.Errorf("expected maxid of 1, got %v", s.MaxId())
	}

	id, err := s.NewTag("foo")
	if err != nil {
		t.Errorf("unexpected error")
	}

	if id != 1 {
		t.Errorf("expected id=1 got %v", id)
	}

	// ensure the id is the same
	id, err = s.NewTag("foo")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}
	if id != 1 {
		t.Errorf("expected id=1 got %v", id)
	}

	// try reinsert
	id, err = s.NewTag("foo")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}
	if id != 1 {
		t.Errorf("expected id=1 got %v", id)
	}

	// insert child node
	id, err = s.NewTag("foo", "bar")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}
	if id != 2 {
		t.Errorf("expected id=2 got %v", id)
	}
}
