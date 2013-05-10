package storetag

import (
	"encoding/gob"
	"io"
)

type Line struct {
	Id     uint64
	Parent uint64
	Name   string
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

	// this skips the "" tag, intentionally
	for i := uint64(1); i < s.maxId; i++ {
		tag := s.ids[i]
		val := Line{tag.id, tag.parent.id, tag.name}

		err := enc.Encode(val)
		if err != nil {
			return err
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
		if err != nil {
			return err
		}

		var tag Tag
		var ok bool
		tag.id = line.Id
		tag.name = line.Name
		tag.parent, ok = s.ids[line.Parent]
		if !ok {
			log.Printf("unknown parent for tag %#v", line)
			continue
		}

		tag.children = make(map[string]*Tag)
		tag.parent.children[line.Name] = &tag

		s.ids[line.Id] = &tag
	}

	return nil
}
