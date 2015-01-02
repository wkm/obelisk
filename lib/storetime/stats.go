package storetime

import (
	"github.com/wkm/obelisk/lib/rinst"
)

var (
	// Stats has measurements on the timestore requests recieved.
	Stats       = rinst.NewCollection()
	statInsert  = Stats.Counter("insert", "points", "datapoints inserted")
	statQuery   = Stats.Counter("query", "op", "query commands received")
	statIter    = Stats.Counter("iter", "points", "points iterated over to fulfill point-level operations (query, dump, etc.)")
	statDump    = Stats.Counter("dump", "op", "dump commands received")
	statLoad    = Stats.Counter("load", "op", "load commands received")
	statError   = Stats.Counter("error", "count", "number of errors")
	statFlush   = Stats.Counter("flush", "op", "flush commands received")
	statCleanup = Stats.Counter("cleanup", "op", "cleanup commands received")
)
