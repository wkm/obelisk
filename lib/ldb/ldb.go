package ldb

import (
	"github.com/jmhodges/levigo"

	"github.com/wkm/obelisk/lib/rinst"
)

// Config contains the exposed configuration options passed to the underlying LevelDB store.
type Config struct {
	Dir       string
	CacheSize int
	Reverse   bool
}

// Store is a wrapper containing instrumentation and handles onto a LevelDB store.
type Store struct {
	cache *levigo.Cache
	DB    *levigo.DB

	Opts struct {
		R      *levigo.ReadOptions
		RCache *levigo.ReadOptions
		WSync  *levigo.WriteOptions
		WAsync *levigo.WriteOptions
	}

	Stats *rinst.Collection

	// Convenience container to allow direct reference to measurement instruments.
	stats struct {
		put   *rinst.Counter
		get   *rinst.Counter
		iter  *rinst.Counter
		write *rinst.Counter
	}
}

// NewStore creates a store with the given configuration.
func NewStore(c Config) (s *Store, err error) {
	s = new(Store)
	s.cache = levigo.NewLRUCache(c.CacheSize)

	opts := levigo.NewOptions()
	opts.SetCache(s.cache)

	opts.SetCreateIfMissing(true) // FIXME silly for this to be datastore level

	s.DB, err = levigo.Open(c.Dir, opts)
	if err != nil {
		s.cache.Close()
		return nil, err
	}

	s.Opts.R = levigo.NewReadOptions()
	s.Opts.R.SetFillCache(false)

	s.Opts.RCache = levigo.NewReadOptions()
	s.Opts.RCache.SetFillCache(true)

	s.Opts.WSync = levigo.NewWriteOptions()
	s.Opts.WSync.SetSync(true)

	s.Opts.WAsync = levigo.NewWriteOptions()
	s.Opts.WAsync.SetSync(false)

	s.Stats = rinst.NewCollection()
	s.stats.get = s.Stats.Counter("get", "op", "Sync and async get requests")
	s.stats.put = s.Stats.Counter("put", "op", "Sync and async put requests")
	s.stats.iter = s.Stats.Counter("iter", "op", "Sync and async iter requests")
	s.stats.write = s.Stats.Counter("write", "op", "Batch write requests")

	// Pull stats from LevelDB using s.DB.PropertyValue("leveldb.stats")

	return
}

// Put synchronously stores an opaque value associated with an opaque key.
func (s *Store) Put(key []byte, value []byte) error {
	s.stats.put.Incr()
	return s.DB.Put(s.Opts.WSync, key, value)
}

// PutAsync stores an opaque value associated with an opaque key.
func (s *Store) PutAsync(key []byte, value []byte) error {
	s.stats.put.Incr()
	return s.DB.Put(s.Opts.WAsync, key, value)
}

// CacheGet gets the value associated with the given key, using the cache if possible.
func (s *Store) CacheGet(key []byte) (value []byte, err error) {
	s.stats.get.Incr()
	return s.DB.Get(s.Opts.RCache, key)
}

// Get the value associated with the given key.
func (s *Store) Get(key []byte) (value []byte, err error) {
	s.stats.get.Incr()
	return s.DB.Get(s.Opts.R, key)
}

// Iterator creates a new levigo.Iterator to enable efficient scans.
func (s *Store) Iterator() (iter *levigo.Iterator) {
	s.stats.iter.Incr()
	return s.DB.NewIterator(s.Opts.R)
}

// CacheIterator creates a new iterator which uses the read cache.
func (s *Store) CacheIterator() (iter *levigo.Iterator) {
	s.stats.iter.Incr()
	return s.DB.NewIterator(s.Opts.RCache)
}

// WriteSync preforms a synchronous batch write.
func (s *Store) WriteSync(wb *levigo.WriteBatch) error {
	s.stats.write.Incr()
	return s.DB.Write(s.Opts.WSync, wb)
}

// WriteAsync preforms a batch write.
func (s *Store) WriteAsync(wb *levigo.WriteBatch) error {
	s.stats.write.Incr()
	return s.DB.Write(s.Opts.WAsync, wb)
}
