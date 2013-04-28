/*
	dead simple in-memory store for tags

	tags are strings with a uniqueid and optionally a parent tag
*/

package storetag

import (
	"errors"
	"sync"
)

type Store struct {
	sync.Mutex
	names map[string]*Tag
	ids   map[uint64]*Tag
	maxId uint64
}

type Tag struct {
	id       uint64
	name     string
	parent   *Tag
	children []*Tag
}

func NewStore() *Store {
	s := new(Store)
	s.names = make(map[string]*Tag)
	s.ids = make(map[uint64]*Tag)
	s.maxId = 0
	return s
}

// give the largest id generated (currently this corresponds to the number of tags)
func (s *Store) MaxId() uint64 {
	s.Lock()
	defer s.Unlock()
	return s.maxId
}

// get the id of a tag (0 on error)
func (s *Store) TagId(name string) (uint64, error) {
	s.Lock()
	defer s.Unlock()
	tag, ok := s.names[name]
	if !ok {
		return 0, errors.New("unknown tag " + name)
	}

	return tag.id, nil
}

// create a new tag with no parent and return its id (0 on error)
func (s *Store) NewTag(name string) (uint64, error) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.names[name]
	if ok {
		return 0, errors.New("node already exists")
	}

	var tag Tag
	s.maxId++
	tag.id = s.maxId
	tag.name = name
	s.names[name] = &tag
	s.ids[tag.id] = &tag

	return tag.id, nil
}

// create a new tag and return its id (0 on error)
func (s *Store) NewChildTag(name string, parent string) (uint64, error) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.names[name]
	if ok {
		return 0, errors.New("node already exists")
	}

	parentTag, ok := s.names[parent]
	if !ok {
		return 0, errors.New("unknown parent " + parent)
	}

	// create a new tag
	var tag Tag
	s.maxId++
	tag.id = s.maxId
	tag.name = name
	tag.parent = parentTag

	// register it with the parent
	parentTag.children = append(parentTag.children, &tag)

	// register with the lookups
	s.names[name] = &tag
	s.ids[tag.id] = &tag

	return tag.id, nil
}

// get the children of a tag
func (s *Store) Children(name string) ([]string, error) {
	node, ok := s.names[name]
	if !ok {
		return nil, errors.New("unknown tag " + name)
	}

	children := make([]string, len(node.children))
	for i, child := range node.children {
		children[i] = child.name
	}

	return children, nil
}
