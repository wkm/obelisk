package storetime

import (
	"circuit/kit/lockfile"
	"log"
	"obelisk/lib/persist"
	"os"
	"time"
)

type Config struct {
	DiskStore     string        // which directory to store the dump files on disk
	FlushPeriod   time.Duration // how often to flush to disk
	FlushVersions int           // how many flushed versions to keep
}

func NewConfig() Config {
	var c Config
	c.FlushVersions = 10
	c.FlushPeriod = 1 * time.Minute
	return c
}

// the storetime DB has persistence feature
type DB struct {
	Store       *Store
	Config      Config
	quit        chan bool
	flushTicker *time.Ticker
	lockFile    *lockfile.LockFile
}

// create a new database
func NewDB(config Config) (*DB, error) {
	store := NewStore()

	err := os.MkdirAll(config.DiskStore, 0700)
	if err != nil {
		return nil, err
	}

	lockfile, err := persist.Lock(config.DiskStore, "storetime")
	if err != nil {
		return nil, err
	}

	db := new(DB)
	db.lockFile = lockfile
	db.Config = config
	db.Store = store

	db.quit = make(chan bool)

	// restore the database
	err = db.Restore()
	if err != nil {
		log.Printf("error restoring: %s", err.Error())
	}

	db.flushTicker = time.NewTicker(config.FlushPeriod)

	go db.backgroundWork()

	return db, nil
}

// background worker for the database
func (db *DB) backgroundWork() {
	for {
		select {
		case <-db.quit:
			return

		case <-db.flushTicker.C:
			log.Printf("flushing to disk")
			statFlush.Incr()
			db.Flush()

			log.Printf("cleaning up flush files")
			statCleanup.Incr()
			db.Cleanup()
		}
	}
}

// load all keys from youngest flush. (in addition to any keys already set)
func (db *DB) Restore() error {
	return persist.RestoreSnapshot(db.Store, db.Config.DiskStore, "time")
}

// flush this db to disk
// FIXME need to include a hash+
func (db *DB) Flush() error {
	statFlush.Incr()
	return persist.FlushSnapshot(db.Store, db.Config.DiskStore, "time")
}

// FIXME implement
func (db *DB) Cleanup() error {
	statCleanup.Incr()
	return persist.CleanupSnapshot(db.Config.FlushVersions, db.Config.DiskStore, "time")
}

// shutdown this store
func (db *DB) Shutdown() {
	close(db.quit)
	db.flushTicker.Stop()

	err := db.lockFile.Release()
	if err != nil {
		log.Printf("error shutting down %s", err.Error())
	}
}
