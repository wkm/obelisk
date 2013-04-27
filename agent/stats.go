package agent

import (
	"obelisk/rinst"
)

var Stats = make(rinst.Collection)

var (
	statMeasurements = Stats.Counter("measures")
	// statMemory       = Stats.Allocation("mem")
)
