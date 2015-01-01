package rlog

import (
	"github.com/wkm/obelisk/lib/rinst"
)

// Stats around logging activity.
var (
	Stats     = rinst.NewCollection()
	statPrint = Stats.Counter("print", "op", "print commands received")
	statFlush = Stats.Counter("flush", "op", "print commands received")
	statByte  = Stats.Counter("byte", "byte", "bytes received into log")
)
