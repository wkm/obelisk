package storetag

import (
	"encoding/binary"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmhodges/levigo"
	"github.com/wkm/obelisk/lib/ldb"
	"github.com/wkm/obelisk/lib/rlog"
)

var log = rlog.LogConfig.Logger("storetag")

type Config struct {
	DiskStore string
	CacheSize int
}

var DefaultConfig = Config{
	CacheSize: 1024 * 1024 * 24,
}

type DB struct {
	mx sync.RWMutex

	config Config
	Store  *ldb.Store
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

// Id gives the id of a tag, if it exists
func (db *DB) Id(name ...string) (id uint64, err error) {
	db.mx.RLock()
	defer db.mx.RUnlock()

	statId.Incr()
	path := createPath(name...)

	bb, err := db.Store.CacheGet([]byte(path))
	if err == nil {
		id = binary.LittleEndian.Uint64(bb)
	}

	return
}

// Get the id of a tag, creating it and the hierarchy to it if it doesn't exist
func (db *DB) Tag(id uint64, name ...string) (newId uint64, err error) {
	statNew.Incr()
	path := []byte(createPath(name...))

	db.mx.Lock()
	defer db.mx.Unlock()

	// Test if the tag exists
	bb, err := db.Store.CacheGet(path)
	if bb != nil {
		// Give the existing tag
		return binary.LittleEndian.Uint64(bb), nil
	}

	// Create the tag
	bb = make([]byte, 8)
	binary.LittleEndian.PutUint64(bb, id)
	err = db.Store.PutAsync(path, bb)
	newId = id

	return
}

// Get the entirety of the tree under a tag
func (db *DB) Children(name ...string) (children []string, err error) {
	db.mx.RLock()
	defer db.mx.RUnlock()

	statChildren.Incr()
	path := createPath(name...)

	children = make([]string, 0, 10)

	iter := db.Store.Iterator()
	defer iter.Close()

	for iter.Seek([]byte(path)); iter.Valid(); iter.Next() {
		child := string(iter.Key())
		if strings.HasPrefix(child, path) {
			children = append(children, child)
		} else {
			break
		}
	}

	err = iter.GetError()
	return
}

// Delete removes a node and all of its children
func (db *DB) Delete(name ...string) (children []string, err error) {
	db.mx.Lock()
	defer db.mx.Unlock()

	statDelete.Incr()

	children, err = db.Children(name...)
	if err != nil {
		return
	}

	wb := levigo.NewWriteBatch()
	defer wb.Close()

	for _, c := range children {
		wb.Delete([]byte(c))
	}

	err = db.Store.WriteAsync(wb)

	return
}

func createPath(name ...string) string {
	return filepath.Join(name...)
}
