package rlog

import (
	"obelisk/rinst"
)

// a collection of rinst stats 
var Stats = make(rinst.Collection)

var (
	statPrints  = Stats.Counter("prints")
	statFlushes = Stats.Counter("flushes")
	statBytes   = Stats.Counter("bytes")
)
