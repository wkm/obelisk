package storekv

import (
	"bytes"
	"encoding/gob"

	"github.com/wkm/obelisk/lib/ldb"
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/rlog"
)

var log = rlog.LogConfig.Logger("storekv")

// Config contains configuration settings for the LevelDB backing.
type Config struct {
	DiskStore string
	CacheSize int
}

// DB is a container that provides KV semantics with a LevelDB backing.
type DB struct {
	config Config
	Store  *ldb.Store
	Stats  *rinst.Collection

	statExists, statGet, statSet *rinst.Counter
}

// NewDB creates a new KV datastore with the given configuration.
func NewDB(config Config) (db *DB, err error) {
	db = new(DB)
	db.config = config
	db.Store, err = ldb.NewStore(ldb.Config{Dir: config.DiskStore, CacheSize: config.CacheSize})

	db.Stats = rinst.NewCollection()
	db.statExists = db.Stats.Counter("exists", "op", "exists commands received")
	db.statGet = db.Stats.Counter("get", "op", "get commands received")
	db.statSet = db.Stats.Counter("set", "op", "set commands received")
	return
}

// Shutdown safely closes the datastore
func (db *DB) Shutdown() {
	db.Store.DB.Close()
	db.Store = nil
}

// Get gives the value for a key, giving empty slice if no such key exists
func (db *DB) Get(key string) (b []byte, err error) {
	db.statGet.Incr()
	b, err = db.Store.CacheGet([]byte(key))
	return
}

// MultiGet gives the values of multiple keys
func (db *DB) MultiGet(keys []string) (values [][]byte, err error) {
	values = make([][]byte, len(keys))

	for i, key := range keys {
		values[i], err = db.Store.CacheGet([]byte(key))
		if err != nil {
			return
		}
	}

	return
}

// Set stores a value under a given key
func (db *DB) Set(key string, value []byte) (err error) {
	db.statSet.Incr()
	err = db.Store.PutAsync([]byte(key), value)
	return
}

// SetGob stores the given object under the given key after Gob encoding.
// FIXME I think this should reuse an encoder for performance reasons
func (db *DB) SetGob(key string, obj interface{}) (err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err = enc.Encode(obj)

	if err != nil {
		return err
	}

	return db.Set(key, b.Bytes())
}

// GetGob gets value of a key and stores it into obj if it's appropriate.
func (db *DB) GetGob(name string, obj interface{}) (err error) {
	value, err := db.Get(name)
	if err != nil {
		return err
	}

	var b = bytes.NewBuffer(value)
	dec := gob.NewDecoder(b)
	return dec.Decode(obj)
}
