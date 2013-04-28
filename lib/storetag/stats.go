package storetag

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statTags    = Stats.Counter("tag")
	statReads   = Stats.Counter("read")
	statWrites  = Stats.Counter("write")
	statFlush   = Stats.Counter("flush")
	statCleanup = Stats.Counter("cleanup")
	statDump    = Stats.Counter("dump")
	statLoad    = Stats.Counter("load")
)