package server

import (
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/rinst/runtime"
	"github.com/wkm/obelisk/lib/storekv"
	"github.com/wkm/obelisk/lib/storetag"
	"github.com/wkm/obelisk/lib/storetime"
)

var (
	// Stats contains measurements from the underlying services of obelisk server.
	Stats = rinst.NewCollection()

	// CommandStats contains measurements of the commands executed.
	CommandStats = rinst.NewCollection()

	// Stats on command requests
	statDeclare = CommandStats.Counter("declare", "op", "declare requests received")
	statSchema  = CommandStats.Counter("schema", "op", "schema requests received")
	statRecord  = CommandStats.Counter("record", "op", "record requests received")
)

func init() {
	Stats.AddInstrument("runtime", runtime.Stats)

	// Include stats from dependencies
	Stats.AddInstrument("kv", storekv.Stats)
	Stats.AddInstrument("tag", storetag.Stats)
	Stats.AddInstrument("time", storetime.Stats)
}
