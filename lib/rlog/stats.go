package rlog

import (
	"obelisk/lib/rinst"
)

// a collection of rinst stats
var Stats = rinst.NewCollection()

var (
	statPrint = Stats.Counter("print", "op", "print commands received")
	statFlush = Stats.Counter("flush", "op", "print commands received")
	statByte  = Stats.Counter("byte", "byte", "bytes received into log")
)
