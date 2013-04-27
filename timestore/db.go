package timestore

import (
	"circuit/kit/lockfile"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// the timestore DB has persistence feature
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

	// create lockfile
	lock, err := lockfile.Create(filepath.Join(config.DiskStore, "lock"))
	if err != nil {
		return nil, errors.New("could not create lock ")
	}

	db := new(DB)
	db.lockFile = lock
	db.config = config
	db.Store = NewStore()

	db.quit = make(chan bool)

	// restore the database
	db.Restore()

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
	matches, err := filepath.Glob(filepath.Join(db.config.DiskStore, "flush-*"))
	if err != nil {
		return err
	}

	if len(matches) < 1 {
		return errors.New("no flushes to restore")
	}

	sort.Strings(matches)

	for i := len(matches) - 1; i >= 0; i-- {
		restoreFile := matches[i]
		log.Printf("attempting restore from %s", restoreFile)

		f, err := os.Open(restoreFile)
		if err != nil {
			log.Printf("  err: %s", err)
			continue
		}
		defer f.Close()

		err = db.Store.Load(f)
		if err != nil {
			log.Printf("  err: %s", err)
			continue
		}

		log.Printf("  restored")
		return nil
	}

	return errors.New("could not successfully restore any flush")
}

// flush this db to disk
// FIXME need to include a hash+
func (db *DB) Flush() {
	statFlush.Incr()

	ts := time.Now().Format(time.RFC3339)
	fname := filepath.Join(db.config.DiskStore, "flush-"+ts)
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
