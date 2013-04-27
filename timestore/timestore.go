/*
	dead simple in-memory store for timeseries, thread safe with a
	global lock; flushes to disk on a regular basis
*/

package timestore

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"io"
	"sync"
)

// a store of named timeseries
type Store struct {
	sync.Mutex
	values map[uint64]*llrb.Tree
}

// a <time,value> pair
type Point struct {
	Time  uint64
	Value float64
}

// create a new timeseries store
func NewStore() *Store {
	s := new(Store)
	s.values = make(map[uint64]*llrb.Tree, 1000)
	return s
}

// inserts the given time and value under the key
func (s *Store) Insert(key, time uint64, value float64) *Store {
	statInsert.Incr()

	s.Lock()
	defer s.Unlock()

	// create a timeseries if we need to
	if _, ok := s.values[key]; !ok {
		s.values[key] = newTree()
	}

	s.values[key].ReplaceOrInsert(Point{time, value})
	return s
}

func newTree() *llrb.Tree {
	return llrb.New(lessPoint)
}

// return all points from key with time in [start,stop]
func (s *Store) Query(key, start, stop uint64) ([]Point, error) {
	statQuery.Incr()

	s.Lock()
	defer s.Unlock()

	var ary []Point
	iterfn := func(i llrb.Item) bool {
		statIter.Incr()

		p := i.(Point)
		if p.Time <= stop {
			ary = append(ary, p)
			return true
		}
		return false
	}

	if ts, ok := s.values[key]; ok {
		ts.AscendGreaterOrEqual(Point{start, 0}, iterfn)
	}

	return ary, nil
}

// return all values from key with time [start,stop]
func (s *Store) FlatQuery(key, start, stop uint64) ([]float64, error) {
	statQuery.Incr()

	s.Lock()
	defer s.Unlock()

	var ary []float64
	iterfn := func(i llrb.Item) bool {
		statIter.Incr()

		p := i.(Point)
		if p.Time <= stop {
			ary = append(ary, p.Value)
			return true
		}
		return false
	}

	if ts, ok := s.values[key]; ok {
		ts.AscendGreaterOrEqual(Point{start, 0}, iterfn)
	}

	return ary, nil
}

// write out all values into a writer, this blocks the entire time
func (s *Store) Dump(w io.Writer) error {
	statDump.Incr()

	s.Lock()
	defer s.Unlock()

	var cnt int
	var err error
	hasErrored := false
	iterfn := func(i llrb.Item) bool {
		// don't iterate if we've errored
		if hasErrored {
			return false
		}

		p := i.(Point)
		cnt, err = fmt.Fprintf(w, "%d=%f,", p.Time, p.Value)
		statIter.Add(uint(cnt))
		if err != nil {
			statError.Incr()
			hasErrored = true
		}

		return true
	}

	// iterate over all keys
	for key, series := range s.values {
		count, err := fmt.Fprintf(w, "%d:", key)
		statDumpBytes.Add(uint(count))
		if err != nil {
			statError.Incr()
			return err
		}

		series.AscendGreaterOrEqual(Point{0, 0}, iterfn)
		if hasErrored {
			return err
		}

		count, err = fmt.Fprintf(w, "\n")
		statDumpBytes.Add(uint(count))
		if err != nil {
			statError.Incr()
			return err
		}
	}

	return nil
}

// sort Point by their time
func lessPoint(a, b interface{}) bool {
	return a.(Point).Time < b.(Point).Time
}
