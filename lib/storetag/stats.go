package storetag

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statId       = Stats.Counter("id", "op", "id commands received")
	statNew      = Stats.Counter("new", "op", "new commands received")
	statChildren = Stats.Counter("children", "op", "children commands received")

	statLoad    = Stats.Counter("load", "op", "load commands received")
	statDump    = Stats.Counter("dump", "op", "dump commands received")
	statFlush   = Stats.Counter("flush", "op", "flush commands received")
	statCleanup = Stats.Counter("cleanup", "op", "cleanup commands received")
)
