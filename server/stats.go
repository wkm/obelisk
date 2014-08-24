package server

import (
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/rinst/runtime"
	"github.com/wkm/obelisk/lib/storekv"
	"github.com/wkm/obelisk/lib/storetag"
	"github.com/wkm/obelisk/lib/storetime"
)

var Stats = rinst.NewCollection()
var CommandStats = rinst.NewCollection()

var (
	StatDeclare = CommandStats.Counter("declare", "op", "declare requests received")
	StatSchema  = CommandStats.Counter("schema", "op", "schema requests received")
	StatRecord  = CommandStats.Counter("record", "op", "record requests received")
)

func init() {
	Stats.AddInstrument("runtime", runtime.Stats)

	// Include stats from dependencies
	Stats.AddInstrument("kv", storekv.Stats)
	Stats.AddInstrument("tag", storetag.Stats)
	Stats.AddInstrument("time", storetime.Stats)

	//
}
