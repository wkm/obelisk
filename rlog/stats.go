package rlog

import (
	"obelisk/rinst"
)

// a collection of rinst stats
var Stats = make(rinst.Collection)

var (
	statPrint = Stats.Counter("print")
	statFlush = Stats.Counter("flush")
	statByte  = Stats.Counter("byte")
)
