package timestore

import (
	"circuit/kit/lockfile"
	"log"
	"os"
	"path/filepath"
	"time"
)

// the timestore DB has persistence feature
type DB struct {
	Store         *Store
	config        Config
	quit          chan bool
	flushTicker   *time.Ticker
	cleanupTicker *time.Ticker
	lockFile      lockfile.LockFile
}

// create a new database
func NewDB(config Config) (*DB, error) {
	err := ValidateConfig(config)
	if err != nil {
		return nil, err
	}

	db := new(DB)
	db.config = config
	db.Store = NewStore()

	// restore the database
	db.Restore()

	db.flushTicker = time.NewTicker(config.FlushPeriod)
	db.cleanupTicker = time.NewTicker(config.CleanupPeriod)

	go db.backgroundWork()

	return db, nil
}

// load all keys from youngest flush. (in addition to any keys already set)
func (db *DB) Restore() error {
	// FIXME implement
	return nil
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

// flush this db to disk
// FIXME need to include a hash+
func (db *DB) Flush() {
	ts := time.Now().Format(time.ANSIC)
	fname := filepath.Join(db.config.DiskStore, ts)
	log.Printf("creating flush %s", fname)

	f, err := os.Create(fname)
	if err != nil {
		log.Printf("could not flush %s", err.Error())
		return
	}
	defer f.Close()

	err = db.Store.Dump(f)
	if err != nil {
		log.Printf("error flushing %s", err.Error())
	}

	log.Printf("flush finished")
}

func (db *DB) Cleanup() {
	log.Printf("cleaning up")
}

// shutdown this store
func (db *DB) Shutdown() {
	close(db.quit)
	db.flushTicker.Stop()
	db.cleanupTicker.Stop()
}
