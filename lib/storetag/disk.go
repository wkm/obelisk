package storetag

import (
	"encoding/gob"
	"io"
)

type Line struct {
	Id   uint64
	Name string

	Parent, Child uint64
}

func init() {
	gob.Register(&Line{})
}

func (s *Store) Dump(w io.Writer) error {
	statDump.Incr()

	s.Lock()
	defer s.Unlock()

	// write out all tag names and their ids
	enc := gob.NewEncoder(w)
	for _, tag := range s.names {
		val := Line{tag.id, tag.name, 0, 0}
		err := enc.Encode(val)
		if err != nil {
			return err
		}
	}

	// write out all tag parents
	for _, tag := range s.names {
		for _, child := range tag.children {
			val := Line{0, "", tag.id, child.id}
			err := enc.Encode(val)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FIXME this should wipe out existing tags
func (s *Store) Load(r io.Reader) error {
	statLoad.Incr()
	dec := gob.NewDecoder(r)

	s.Lock()
	defer s.Unlock()

	for {
		var line Line
		err := dec.Decode(&line)
		if err == io.EOF {
			break
		}

		if line.Id != 0 {
			// new node
			var tag Tag
			tag.id = line.Id
			tag.name = line.Name

			s.names[tag.name] = &tag
			s.ids[tag.id] = &tag
		} else {
			// new edge
			child := s.ids[line.Child]
			parent := s.ids[line.Parent]

			child.parent = parent
			parent.children = append(parent.children, child)
		}
	}

	return nil
}
