package storetag

import (
	"circuit/kit/lockfile"
	"log"
	"obelisk/lib/persist"
	"os"
	"time"
)

type Config struct {
	DiskStore     string
	FlushPeriod   time.Duration
	FlushVersions int
	CleanupPeriod time.Duration
}

func NewConfig() Config {
	var c Config
	c.FlushPeriod = 1 * time.Minute
	c.FlushVersions = 10
	c.CleanupPeriod = 10 * time.Minute
	return c
}

type DB struct {
	Store         *Store
	Config        Config
	quit          chan bool
	flushTicker   *time.Ticker
	cleanupTicker *time.Ticker
	lockFile      *lockfile.LockFile
}

func NewDB(config Config) (*DB, error) {
	store := NewStore()
	db := new(DB)

	db.Store = store
	db.Config = config

	err := os.MkdirAll(config.DiskStore, 0700)
	if err != nil {
		return nil, err
	}

	lockFile, err := persist.Lock(config.DiskStore, "storetag")
	if err != nil {
		return nil, err
	}
	db.lockFile = lockFile

	db.quit = make(chan bool)

	err = db.Restore()
	if err != nil {
		log.Printf("error restoring %s", err.Error())
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
	return persist.RestoreSnapshot(db.Store, db.Config.DiskStore, "tag")
}

// flush this db to disk
func (db *DB) Flush() error {
	statFlush.Incr()
	return persist.FlushSnapshot(db.Store, db.Config.DiskStore, "tag")
}

func (db *DB) Cleanup() error {
	statCleanup.Incr()
	return persist.CleanupSnapshot(db.Config.DiskStore, "tag")
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
