package server

import (
	"obelisk/lib/rinst"
	"obelisk/lib/storekv"
	"obelisk/lib/storetag"
	"obelisk/lib/storetime"
)

var Stats = rinst.NewCollection()

func init() {
	// include stats from dependencies
	Stats.AddInstrument("kv", storekv.Stats)
	Stats.AddInstrument("tag", storetag.Stats)
	Stats.AddInstrument("time", storetime.Stats)
}
