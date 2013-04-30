/*
	dead simple in-memory store for string-keys and []byte blogs
*/

package storekv

import (
	"bytes"
	"encoding/gob"
	"obelisk/lib/errors"
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

// get the values of multiple names
func (s *Store) MultiGet(names []string) ([][]byte, error) {
	values := make([][]byte, len(names))
	s.Lock()
	s.Unlock()

	for i, name := range names {
		values[i] = s.Values[name]
	}

	return values, nil
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

// FIXME I think this should reuse an encoder for performance reasons
func (s *Store) SetGob(name string, obj interface{}) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(obj)

	if err != nil {
		return errors.W(err)
	}

	return s.Set(name, b.Bytes())
}

// get the gob value of a key into obj
func (s *Store) GetGob(name string, obj interface{}) error {
	value, err := s.Get(name)
	if err != nil {
		return errors.W(err)
	}

	var b = bytes.NewBuffer(value)
	dec := gob.NewDecoder(b)
	return dec.Decode(obj)
}
