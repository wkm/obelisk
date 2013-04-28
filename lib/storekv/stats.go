package storekv

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statKeys    = Stats.IntValue("keys")
	statExists  = Stats.Counter("exists")
	statGet     = Stats.Counter("get")
	statSet     = Stats.Counter("set")
	statLoad    = Stats.Counter("load")    // loads executed
	statDump    = Stats.Counter("dump")    // dumps executed
	statFlush   = Stats.Counter("flush")   // number of flushes
	statCleanup = Stats.Counter("cleanup") // number of cleanups
)
