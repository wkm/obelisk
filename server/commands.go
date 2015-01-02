package server

import (
	"errors"
	"math/rand"
	"path/filepath"

	"github.com/wkm/obelisk/lib/rlog"
)

var (
	DuplicateIdentifier = errors.New("Identifier is not unique on path")
)

var log = rlog.LogConfig.Logger("obelisk-server")

// Associate an identifier with the given paths
func (app *ServerApp) Declare(id string, paths ...string) (err error) {
	statDeclare.Incr()
	uid := uint64(rand.Int63())
	for _, path := range paths {
		actid, err := app.tagdb.Tag(uid, filepath.Join(path, id))
		if uid != actid {
			return DuplicateIdentifier
		}
		if err != nil {
			return err
		}
	}
	return
}

// Schema stores metadata on the a metric's structure.
func (app *App) Schema(id, op, kind, unit, desc string) (err error) {
	statSchema.Incr()
	return
}

// Record stores a single measurement of a metric.
func (app *App) Record(id, metric, time, value string) (err error) {
	statRecord.Incr()
	return
}
