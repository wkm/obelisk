/*
	dead simple in-memory store for string-keys and []byte blogs
*/

package storekv

import (
	"errors"
	"sync"
)

type Store struct {
	sync.Mutex
	Values map[string][]byte
}

// create a new kv store
func NewStore() *Store {
	s := new(Store)
	s.Values = make(map[string][]byte)
	return s
}

// get the value of a name
func (s *Store) Get(name string) ([]byte, error) {
	statGet.Incr()

	s.Lock()
	defer s.Unlock()

	b, ok := s.Values[name]
	if ok {
		return b, nil
	} else {
		return nil, errors.New("unknown key")
	}
}

// gives true if the key exists
func (s *Store) Exists(name string) bool {
	statExists.Incr()

	s.Lock()
	defer s.Unlock()

	_, ok := s.Values[name]
	return ok
}

// set the value with a name
func (s *Store) Set(name string, value []byte) error {
	statSet.Incr()

	s.Lock()
	defer s.Unlock()

	s.Values[name] = value

	return nil
}
