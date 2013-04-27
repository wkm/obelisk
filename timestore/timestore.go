/*
	dead simple in-memory store for timeseries, thread safe with a
	global lock; flushes to disk on a regular basis
*/

package timestore

import (
	"github.com/petar/GoLLRB/llrb"
	"sync"
	"time"
)

// a store of named timeseries
type Store struct {
	sync.Mutex
	config    Config
	values    map[uint64]*llrb.Tree
	shutdown  chan bool
	flushTick <-chan time.Time
}

// a <time,value> pair
type Point struct {
	Time  uint64
	Value float64
}

// create a new timeseries store
func NewStore(config Config) *Store {
	s := new(Store)
	s.config = config
	s.values = make(map[uint64]*llrb.Tree, config.InitialSize)

	// load from disk if applicable

	s.flushTick = time.Tick(config.FlushPeriod)
	go s.startPeriodicFlush()
	return s
}

// shutdown this store
func (s *Store) Shutdown() {
	s.shutdown <- true
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

// sort Point by their time
func lessPoint(a, b interface{}) bool {
	return a.(Point).Time < b.(Point).Time
}
