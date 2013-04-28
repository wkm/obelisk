package persist

import (
	"circuit/kit/lockfile"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Persistable interface {
	// dump the current state into the given writer
	Dump(w io.Writer) error

	// load the current state from the given reader
	Load(r io.Reader) error
}

// create a lockfile
func Lock(dir, key string) (*lockfile.LockFile, error) {
	return lockfile.Create(filepath.Join(dir, key))
}

func CleanupSnapshot(dir, key string) error {
	// FIXME implement
	return nil
}

func FlushSnapshot(p Persistable, dir, key string) error {
	ts := time.Now().Format(time.RFC3339)
	fname := filepath.Join(dir, key+"-"+ts)
	log.Printf("creating flush %s", fname)

	f, err := os.Create(fname)
	if err != nil {
		log.Printf("could not flush %s", err.Error())
		return err
	}
	defer f.Close()

	err = p.Dump(f)
	if err != nil {
		log.Printf("error flushing %s", err.Error())
		return err
	}

	log.Printf("flushed")
	return nil
}

func RestoreSnapshot(p Persistable, dir, key string) error {
	searchpath := filepath.Join(dir, key+"-*")

	matches, err := filepath.Glob(searchpath)
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

		err = p.Load(f)
		if err != nil {
			log.Printf("  err: %s", err)
			continue
		}

		log.Printf("  restored")
		return nil
	}

	return errors.New("could not successfully restore any flush")
}