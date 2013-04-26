/*
	dead simple in-memory store for timeseries, thread safe with a
	global lock
*/

package timestore

import (
	"github.com/petar/GoLLRB/llrb"
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
	s.Lock()
	defer s.Unlock()

	// create a timeseries if we need to
	if _, ok := s.values[key]; ok {
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
	s.Lock()
	defer s.Unlock()

	var ary []Point
	iterfn := func(i llrb.Item) bool {
		p := i.(Point)
		if p.Time <= stop {
			ary = append(ary, p)
			return true
		}
		return false
	}

	if ts, ok := s.values[key]; ok {
		ts.AscendGreaterOrEqual(start, iterfn)
	}

	return ary, nil
}

// return all values from key with time [start,stop]
func (s *Store) FlatQuery(key, start, stop uint64) ([]float64, error) {
	s.Lock()
	defer s.Unlock()

	var ary []float64
	iterfn := func(i llrb.Item) bool {
		p := i.(Point)
		if p.Time <= stop {
			ary = append(ary, p.Value)
			return true
		}
		return false
	}

	if ts, ok := s.values[key]; ok {
		ts.AscendGreaterOrEqual(start, iterfn)
	}

	return ary, nil
}

// sorts Point by their time
func lessPoint(a, b interface{}) bool {
	return a.(Point).Time < b.(Point).Time
}
