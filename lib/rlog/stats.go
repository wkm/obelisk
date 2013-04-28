package rlog

import (
	"obelisk/lib/rinst"
)

// a collection of rinst stats
var Stats = rinst.NewCollection()

var (
	statPrint = Stats.Counter("print")
	statFlush = Stats.Counter("flush")
	statByte  = Stats.Counter("byte")
)
