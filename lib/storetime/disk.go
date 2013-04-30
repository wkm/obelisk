package storetime

import (
	"encoding/gob"
	"github.com/petar/GoLLRB/llrb"
	"io"
)

// a fully qualified point which
type FullPoint struct {
	Key, Time uint64
	Value     float64
}

func init() {
	gob.Register(&FullPoint{})
}

// write out all values into a writer, this blocks the entire time
func (s *Store) Dump(w io.Writer) error {
	statDump.Incr()

	s.Lock()
	defer s.Unlock()

	enc := gob.NewEncoder(w)

	var err error
	hasErrored := false

	// iterate over all keys
	for key, series := range s.values {
		iterfn := func(i llrb.Item) bool {
			// don't iterate if we've errored
			if hasErrored {
				return false
			}

			p := i.(Point)
			full := FullPoint{key, p.Time, p.Value}
			err := enc.Encode(full)
			if err != nil {
				statError.Incr()
				hasErrored = true
			}

			return true
		}

		if err != nil {
			statError.Incr()
			return err
		}

		series.AscendGreaterOrEqual(Point{0, 0}, iterfn)
		if hasErrored {
			return err
		}

		if err != nil {
			statError.Incr()
			return err
		}
	}

	return nil
}

// load up values from a dump
func (s *Store) Load(r io.Reader) error {
	// we don't need to lock because our individual inserts will
	statLoad.Incr()
	dec := gob.NewDecoder(r)

	for {
		var point FullPoint
		err := dec.Decode(&point)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}

		s.Insert(point.Key, point.Time, point.Value)
	}

	return nil
}
