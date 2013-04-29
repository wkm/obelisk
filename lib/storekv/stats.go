package storekv

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statExists  = Stats.Counter("exists", "op", "exists commands received")
	statGet     = Stats.Counter("get", "op", "get commands received")
	statSet     = Stats.Counter("set", "op", "set commands received")
	statLoad    = Stats.Counter("load", "op", "load commands received")
	statDump    = Stats.Counter("dump", "op", "dump commands received")
	statFlush   = Stats.Counter("flush", "op", "flush commands received")
	statCleanup = Stats.Counter("cleanup", "op", "cleanup commands received")
)
