package storetime

import (
	"circuit/kit/lockfile"
	"log"
	"obelisk/lib/persist"
	"time"
)

// the storetime DB has persistence feature
type DB struct {
	Store         *Store
	config        Config
	quit          chan bool
	flushTicker   *time.Ticker
	cleanupTicker *time.Ticker
	lockFile      *lockfile.LockFile
}

// create a new database
func NewDB(config Config) (*DB, error) {
	err := ValidateConfig(config)
	if err != nil {
		return nil, err
	}

	store := NewStore()
	lockfile, err := persist.Lock(config.DiskStore, "storetime")
	if err != nil {
		return nil, err
	}

	db := new(DB)
	db.lockFile = lockfile
	db.config = config
	db.Store = store

	db.quit = make(chan bool)

	// restore the database
	err = db.Restore()
	if err != nil {
		log.Printf("error restoring: %s", err.Error())
	}

	db.flushTicker = time.NewTicker(config.FlushPeriod)
	db.cleanupTicker = time.NewTicker(config.CleanupPeriod)

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

		case <-db.cleanupTicker.C:
			log.Printf("cleaning up flush files")
			statCleanup.Incr()
			db.Cleanup()
		}
	}
}

// load all keys from youngest flush. (in addition to any keys already set)
func (db *DB) Restore() error {
	return persist.RestoreSnapshot(db.Store, db.config.DiskStore, "time")
}

// flush this db to disk
// FIXME need to include a hash+
func (db *DB) Flush() error {
	statFlush.Incr()
	return persist.FlushSnapshot(db.Store, db.config.DiskStore, "time")
}

// FIXME implement
func (db *DB) Cleanup() {
	log.Printf("cleaning up")
}

// shutdown this store
func (db *DB) Shutdown() {
	close(db.quit)
	db.flushTicker.Stop()
	db.cleanupTicker.Stop()

	err := db.lockFile.Release()
	if err != nil {
		log.Printf("error shutting down %s", err.Error())
	}
}
