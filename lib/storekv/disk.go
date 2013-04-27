package storekv

import (
	"encoding/gob"
	"io"
)

type FullRow struct {
	Key   string
	Value []byte
}

func init() {
	gob.Register(&FullRow{})
}

func (s *Store) Dump(w io.Writer) error {
	statDump.Incr()

	s.Lock()
	defer s.Unlock()

	enc := gob.NewEncoder(w)

	for k, v := range s.Values {
		err := enc.Encode(FullRow{k, v})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Load(r io.Reader) error {
	statLoad.Incr()
	dec := gob.NewDecoder(r)

	for {
		var row FullRow
		err := dec.Decode(&row)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		s.Set(row.Key, row.Value)
	}

	return nil
}
