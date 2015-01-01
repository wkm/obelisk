package ldb

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	cacheSize = 1024 * 1024 * 24
)

func TestStore(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "levelkv")
	defer os.RemoveAll(dir)

	s, err := NewStore(Config{Dir: dir, CacheSize: cacheSize})
	if err != nil {
		t.Fatal(err)
	}

	if err = s.Put([]byte("key"), []byte("value")); err != nil {
		t.Fatal(err)
	}

	if err = s.Put([]byte("key"), []byte("value2")); err != nil {
		t.Fatal(err)
	}

	res, err := s.Get([]byte("key"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("got %s", res)
	if string(res) != "value2" {
		t.Errorf("expected %#v got %s", "value2", res)
	}
}

// A trivial benchmark to establish some kind of baseline performance
func BenchmarkSyncPut(b *testing.B) {
	dir, err := ioutil.TempDir(os.TempDir(), "levelkv")
	s, err := NewStore(Config{Dir: dir, CacheSize: cacheSize})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Put([]byte("we've got a key. yay"), []byte("this is a medium sized value"))
	}
}

func BenchmarkAsyncPut(b *testing.B) {
	dir, err := ioutil.TempDir(os.TempDir(), "levelkv")
	s, err := NewStore(Config{Dir: dir, CacheSize: cacheSize})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.PutAsync([]byte("we've got a key. yay"), []byte("this is a medium sized value"))
	}
}

func BenchmarkGet(b *testing.B) {
	dir, err := ioutil.TempDir(os.TempDir(), "levelkv")
	s, err := NewStore(Config{Dir: dir, CacheSize: cacheSize})
	if err != nil {
		b.Fatal(err)
	}

	s.PutAsync([]byte("we've got a key. yay"), []byte("this is a medium sized value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get([]byte("we've got a key. yay"))
	}
}

func BenchmarkCachedGet(b *testing.B) {
	dir, err := ioutil.TempDir(os.TempDir(), "levelkv")
	s, err := NewStore(Config{Dir: dir, CacheSize: cacheSize})
	if err != nil {
		b.Fatal(err)
	}

	s.PutAsync([]byte("we've got a key. yay"), []byte("this is a medium sized value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.CacheGet([]byte("we've got a key. yay"))
	}
}
