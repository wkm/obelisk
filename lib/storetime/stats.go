package storetime

import (
	"obelisk/lib/rinst"
)

var Stats = make(rinst.Collection)

var (
	statInsert  = Stats.Counter("insert")  // datapoints inserted
	statQuery   = Stats.Counter("query")   // queries executed
	statIter    = Stats.Counter("iter")    // points iterated over to fulfill queries
	statDump    = Stats.Counter("dump")    // dumps executed
	statLoad    = Stats.Counter("load")    // loads executed
	statError   = Stats.Counter("error")   // number of errors encountered
	statFlush   = Stats.Counter("flush")   // number of flushes
	statCleanup = Stats.Counter("cleanup") // number of cleanups
)
