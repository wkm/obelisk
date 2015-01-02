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
	StatDeclare.Incr()
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

func (app *ServerApp) Schema(id, op, kind, unit, desc string) (err error) {
	StatSchema.Incr()
	return
}

func (app *ServerApp) Record(id, metric, time, value string) (err error) {
	StatRecord.Incr()
	return
}
