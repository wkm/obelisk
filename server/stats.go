package server

import (
	"github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/rinst/runtime"
	"github.com/wkm/obelisk/lib/storekv"
	"github.com/wkm/obelisk/lib/storetag"
	"github.com/wkm/obelisk/lib/storetime"
)

var Stats = rinst.NewCollection()

func init() {
	Stats.AddInstrument("runtime", runtime.Stats)

	// include stats from dependencies
	Stats.AddInstrument("kv", storekv.Stats)
	Stats.AddInstrument("tag", storetag.Stats)
	Stats.AddInstrument("time", storetime.Stats)
}
