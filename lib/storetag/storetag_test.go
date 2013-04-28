package storetag

import (
	"testing"
)

func TestTagStore(t *testing.T) {
	s := NewStore()

	if s.MaxId() != 0 {
		t.Errorf("expected maxid of 0, got %v", s.MaxId())
	}

	id, err := s.NewTag("foo")
	if err != nil {
		t.Errorf("unexpected error")
	}

	if id != 1 {
		t.Errorf("expected id=1 got %v", id)
	}

	// try retrieval
	id, err = s.TagId("foo")
	if err != nil {
		t.Errorf("unexpected error")
	}
	if id != 1 {
		t.Errorf("expected id=1 got %v", id)
	}

	// try reinsert
	id, err = s.NewTag("foo")
	if err == nil {
		t.Errorf("expected an error")
	}
	if id != 0 {
		t.Errorf("expected id=0 got %v", id)
	}

	// insert child node
	id, err = s.NewChildTag("bar", "foo")
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}
	if id != 2 {
		t.Errorf("expected id=2 got %v", id)
	}
}
