package server

import (
	"github.com/wkm/obelisk/lib/rlog"
)

var log = rlog.LogConfig.Logger("obelisk-server")

// Associate an identifier with the given paths
func (app *ServerApp) Declare(id string, paths ...string) (err error) {
	StatDeclare.Incr()
	return
	// for _, path := range paths {
	// 	actid, err := app.tagdb.Tag(id, path)
	// 	if id != actid {
	// 		return errors.New("identifier is not unique within path")
	// 	}
	// }
}

func (app *ServerApp) Schema(id, op, kind, unit, desc string) (err error) {
	StatSchema.Incr()
	return
}

func (app *ServerApp) Record(id, metric, time, value string) (err error) {
	StatRecord.Incr()
	return
}
