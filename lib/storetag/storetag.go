/*
	dead simple in-memory store for tags

	tags are used to encode hierarchies, ensuring that every
	node has a unique, persistent ID
*/

package storetag

import (
	"obelisk/lib/errors"
	"obelisk/lib/rlog"
	"path/filepath"
	"strings"
	"sync"
)

var log = rlog.LogConfig.Logger("storetag")

type Store struct {
	sync.Mutex
	ids   map[uint64]*Tag
	maxId uint64

	root *Tag
}

type Tag struct {
	id       uint64
	name     string
	parent   *Tag
	children map[string]*Tag
}

func NewStore() *Store {
	s := new(Store)
	s.ids = make(map[uint64]*Tag)
	s.maxId = 0

	// create the root node
	s.root = s.newTag("")

	return s
}

// give the largest id generated (one less than the number of tags)
func (s *Store) MaxId() uint64 {
	s.Lock()
	defer s.Unlock()

	return s.maxId
}

// not threadsafe
func (s *Store) newTag(name string) *Tag {
	var tag Tag
	tag.id = s.maxId
	s.maxId++
	s.ids[tag.id] = &tag

	tag.name = name
	tag.children = make(map[string]*Tag)

	return &tag
}

func createPath(name ...string) []string {
	return strings.Split("/"+filepath.Join(name...), "/")
}

// get the id of a tag, if it exists
func (s *Store) Id(name ...string) (uint64, error) {
	statId.Incr()
	components := createPath(name...)

	s.Lock()
	defer s.Unlock()

	cursor := s.root
	for _, part := range components[1:] {
		child, ok := cursor.children[part]
		if !ok {
			return 0, errors.N("unknown node " + part + " of " + strings.Join(components, "/"))
		}
		cursor = child
	}

	return cursor.id, nil
}

// get the id of a tag, creating it and the hierarchy to if it doesn't exist
func (s *Store) NewTag(name ...string) (uint64, error) {
	statNew.Incr()
	components := createPath(name...)

	s.Lock()
	defer s.Unlock()

	cursor := s.root
	for _, part := range components[1:] {
		child, ok := cursor.children[part]
		if !ok {
			tag := s.newTag(part)
			tag.parent = cursor
			cursor.children[part] = tag
			s.ids[tag.id] = tag
			child = tag
		}

		cursor = child
	}

	return cursor.id, nil
}

// get the children of a tag
func (s *Store) Children(name ...string) ([]string, error) {
	statChildren.Incr()
	components := createPath(name...)

	s.Lock()
	defer s.Unlock()

	cursor := s.root
	for _, part := range components[1:] {
		child, ok := cursor.children[part]
		if !ok {
			return nil, errors.N("unknown node " + part + " of " + strings.Join(components, "/"))
		}

		cursor = child
	}

	childrenNames := make([]string, len(cursor.children))
	i := 0
	for _, c := range cursor.children {
		childrenNames[i] = c.name
		i++
	}

	return childrenNames, nil
}

// delete a node
func (s *Store) Delete(name ...string) error {
	statDelete.Incr()
	components := createPath(name...)

	s.Lock()
	defer s.Unlock()

	cursor := s.root
	for _, part := range components[1:] {
		child, ok := cursor.children[part]
		if !ok {
			return errors.N("unknown node " + part + " of " + strings.Join(components, "/"))
		}

		cursor = child
	}

	delete(cursor.parent.children, cursor.name)
	cursor.parent = nil

	return nil
}
