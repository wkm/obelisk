package storetime

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/wkm/obelisk/lib/ldb"
	"github.com/wkm/obelisk/lib/rlog"
)

var log = rlog.LogConfig.Logger("storetime")

type Config struct {
	DiskStore string
	CacheSize int
}

type DB struct {
	config Config
	Store  *ldb.Store
}

type Point struct {
	Time  uint64
	Value float64
}

func NewDB(config Config) (db *DB, err error) {
	db = new(DB)
	db.config = config
	db.Store, err = ldb.NewStore(ldb.Config{Dir: config.DiskStore, CacheSize: config.CacheSize})
	return
}

// Safely close the datastore
func (db *DB) Shutdown() {
	db.Store.DB.Close()
	db.Store = nil
}

func createKey(id, time int64) []byte {
	return []byte(fmt.Sprintf("%d•%d", id, time))
}

func getTime(key []byte) uint64 {
	var id, time int64
	fmt.Sscanf(string(key), "%d•%d", &id, &time)
	return time
}

func (db *DB) Insert(key, time int64, value float64) {
	statInsert.Incr()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(value))
	db.Store.PutAsync(createKey(key, time), b)
}

// Query gives all tuples of <time,value> from key with time in [start,stop]
func (db *DB) Query(key, start, stop uint64) (points []Point, err error) {
	statQuery.Incr()

	iter := db.Store.CacheIterator()
	end := []byte(createKey(key, stop))

	points = make([]Point, 0, 10)
	for iter.Seek(createKey(key, start)); iter.Valid(); iter.Next() {
		statIter.Incr()
		if bytes.Compare(iter.Key(), end) > 0 {
			break
		}

		p := Point{
			Time:  getTime(iter.Key()),
			Value: math.Float64frombits(binary.LittleEndian.Uint64(iter.Value())),
		}
		points = append(points, p)
	}

	return
}

// FlatQuery gives all values from key with time in [start, stop]
func (db *DB) FlatQuery(key, start, stop uint64) (values []float64, err error) {
	statQuery.Incr()

	iter := db.Store.CacheIterator()
	end := []byte(createKey(key, stop))

	values = make([]float64, 0, 10)
	for iter.Seek(createKey(key, start)); iter.Valid(); iter.Next() {
		statIter.Incr()
		if bytes.Compare(iter.Key(), end) > 0 {
			break
		}

		p := math.Float64frombits(binary.LittleEndian.Uint64(iter.Value()))
		values = append(values, p)
	}

	return
}
