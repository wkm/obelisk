package agent

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statMeasurements = Stats.Counter("measures")
	// statMemory       = Stats.Allocation("mem")
)
