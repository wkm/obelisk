package server

import (
	"obelisk/lib/rinst"
	"obelisk/lib/rinst/runtime"
	"obelisk/lib/storekv"
	"obelisk/lib/storetag"
	"obelisk/lib/storetime"
)

var Stats = rinst.NewCollection()

func init() {
	Stats.AddInstrument("runtime", runtime.Stats)

	// include stats from dependencies
	Stats.AddInstrument("kv", storekv.Stats)
	Stats.AddInstrument("tag", storetag.Stats)
	Stats.AddInstrument("time", storetime.Stats)
}
