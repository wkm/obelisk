package ldb

import (
	"github.com/jmhodges/levigo"
	"github.com/wkm/obelisk/lib/rinst"
)

type Config struct {
	Dir       string
	CacheSize int
	Reverse   bool
}

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

func (s *Store) Put(key []byte, value []byte) error {
	s.stats.put.Incr()
	return s.DB.Put(s.Opts.WSync, key, value)
}

func (s *Store) PutAsync(key []byte, value []byte) error {
	s.stats.put.Incr()
	return s.DB.Put(s.Opts.WAsync, key, value)
}

func (s *Store) CacheGet(key []byte) (value []byte, err error) {
	s.stats.get.Incr()
	return s.DB.Get(s.Opts.RCache, key)
}

func (s *Store) Get(key []byte) (value []byte, err error) {
	s.stats.get.Incr()
	return s.DB.Get(s.Opts.R, key)
}

func (s *Store) Iterator() (iter *levigo.Iterator) {
	s.stats.iter.Incr()
	return s.DB.NewIterator(s.Opts.R)
}

func (s *Store) CacheIterator() (iter *levigo.Iterator) {
	s.stats.iter.Incr()
	return s.DB.NewIterator(s.Opts.RCache)
}

func (s *Store) WriteSync(wb *levigo.WriteBatch) error {
	s.stats.write.Incr()
	return s.DB.Write(s.Opts.WSync, wb)
}

func (s *Store) WriteAsync(wb *levigo.WriteBatch) error {
	s.stats.write.Incr()
	return s.DB.Write(s.Opts.WAsync, wb)
}
