package server

import (
	"errors"
	"math/rand"
	"path/filepath"

	"github.com/wkm/obelisk/lib/rlog"
)

var (
	// ErrDuplicateIdentifier is returned if the specified identifier was already used.
	ErrDuplicateIdentifier = errors.New("Identifier is not unique on path")
)

var log = rlog.LogConfig.Logger("obelisk-server")

// Declare associates an identifier with the given paths.
func (app *App) Declare(id string, paths ...string) (err error) {
	statDeclare.Incr()
	uid := uint64(rand.Int63())
	for _, path := range paths {
		actid, err := app.tagdb.Tag(uid, filepath.Join(path, id))
		if uid != actid {
			return ErrDuplicateIdentifier
		}
		if err != nil {
			return err
		}
	}
	return
}

// Schema stores metadata on the a metric's structure.
func (app *App) Schema(id, kind, unit, desc string) (err error) {
	statSchema.Incr()
	return
}

// Record stores a single measurement of a metric.
func (app *App) Record(id, metric, time, value string) (err error) {
	statRecord.Incr()
	return
}

// KVGet gives the value of key if
func (app *App) KVGet(key string) (str string, err error) {
	bb, err := app.kvdb.Get(key)
	str = string(bb)
	return
}

func (app *App) KVSet(key, value string) (err error) {
	return app.kvdb.Set(key, []byte(value))
}
