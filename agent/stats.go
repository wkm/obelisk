package agent

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statMeasurements = Stats.Counter("measures", "meas", "number of agent measurements")
	// statMemory       = Stats.Allocation("mem")
)
