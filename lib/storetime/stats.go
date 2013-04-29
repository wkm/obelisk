package storetime

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statInsert  = Stats.Counter("insert", "points", "datapoints inserted")
	statQuery   = Stats.Counter("query", "op", "query commands received")
	statIter    = Stats.Counter("iter", "points", "points iterated over to fulfill ")
	statDump    = Stats.Counter("dump", "op", "dump commands received")
	statLoad    = Stats.Counter("load", "op", "load commands received")
	statError   = Stats.Counter("error", "count", "number of errors")
	statFlush   = Stats.Counter("flush", "op", "flush commands received")
	statCleanup = Stats.Counter("cleanup", "op", "cleanup commands received")
)
