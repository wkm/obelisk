package server

import (
	"errors"
	"fmt"
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
func (app *App) Declare(id string, paths []string) (err error) {
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

	app.kvdb.Set(fmt.Sprintf("id•%s"), []byte(fmt.Sprintf("%d", uid)))

	return
}

// Schema stores metadata on the a metric's structure.
func (app *App) Schema(id, kind, unit, desc string) (err error) {
	statSchema.Incr()
	return
}

// Record stores a single measurement of a metric.
func (app *App) Record(id, time, value string) (err error) {
	statRecord.Incr()
	return
}

// KVGet gives the value stored with key.
func (app *App) KVGet(key string) (str string) {
	bb, _ := app.kvdb.Get(key)
	str = string(bb)
	return
}

// KVSet sets the value stored with key.
func (app *App) KVSet(key, value string) (err error) {
	return app.kvdb.Set(key, []byte(value))
}

// TagID gives the identifier for the name.
func (app *App) TagID(name string) (id uint64) {
	id, _ = app.tagdb.ID(name)
	return
}

// TagChildren gives all children of the given tag, in lexicographic order.
func (app *App) TagChildren(path string) []string {
	children, _ := app.tagdb.Children(path)
	return children
}
