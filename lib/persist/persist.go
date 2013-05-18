package persist

import (
	"circuit/kit/lockfile"
	"compress/gzip"
	"io"
	"obelisk/lib/errors"
	"obelisk/lib/rlog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var log = rlog.LogConfig.Logger("persist")

type Persistable interface {
	// dump the current state into the given writer
	Dump(w io.Writer) error

	// load the current state from the given reader
	Load(r io.Reader) error
}

// create a lockfile
func Lock(dir, key string) (*lockfile.LockFile, error) {
	return lockfile.Create(filepath.Join(dir, key+".lock"))
}

func CleanupSnapshot(flushes int, dir, key string) error {
	searchpath := filepath.Join(dir, key+"-*")
	matches, err := filepath.Glob(searchpath)
	if err != nil {
		return err
	}

	if len(matches) < flushes {
		// nothing to cleanup
		return nil
	}

	sort.Strings(matches)
	for _, path := range matches[:len(matches)-flushes] {
		err := os.Remove(path)
		if err != nil {
			return err
		}
		log.Printf("cleaned up %s", path)
	}

	return nil
}

// dump a persistable object into a snapshot on disk
func FlushSnapshot(p Persistable, dir, key string) error {
	ts := time.Now().Format(time.RFC3339)
	fname := filepath.Join(dir, key+"-"+ts+".gz")

	// write to a temp file initially
	tmp := fname + ".tmp"
	log.Printf("creating flush %s", fname)

	f, err := os.Create(tmp)
	if err != nil {
		log.Printf("could not flush %s", err.Error())
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	err = p.Dump(gz)
	if err != nil {
		log.Printf("error flushing %s", err.Error())
		return err
	}

	// save this flush
	os.Rename(tmp, fname)
	log.Printf("flushed")
	return nil
}

// read the last available snapshot from a persistence directory
func RestoreSnapshot(p Persistable, dir, key string) error {
	searchpath := filepath.Join(dir, key+"-*")

	matches, err := filepath.Glob(searchpath)
	if err != nil {
		return err
	}

	if len(matches) < 1 {
		return errors.N("no flushes to restore")
	}

	sort.Strings(matches)

	for i := len(matches) - 1; i >= 0; i-- {
		restoreFile := matches[i]

		// skip temporary files
		if strings.HasSuffix(restoreFile, ".tmp") {
			continue
		}

		log.Printf("attempting restore from %s", restoreFile)

		var r io.Reader
		f, err := os.Open(restoreFile)
		if err != nil {
			log.Printf("err: %s", err)
			continue
		}
		defer f.Close()

		// are we reading a gzip?
		if strings.HasSuffix(restoreFile, ".gz") {
			gz, err := gzip.NewReader(f)
			if err != nil {
				log.Printf("gz-err: %s", err)
				continue
			}
			defer gz.Close()
			r = gz
		} else {
			// read from the file directly
			r = f
		}

		err = p.Load(r)
		if err != nil {
			log.Printf("err: %s", err)
			continue
		}

		log.Printf("restored")
		return nil
	}

	return errors.N("could not successfully restore any flush")
}
