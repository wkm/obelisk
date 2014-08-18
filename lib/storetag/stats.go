package storetag

import (
	"github.com/wkm/obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statId       = Stats.Counter("id", "op", "id commands received")
	statNew      = Stats.Counter("new", "op", "new commands received")
	statChildren = Stats.Counter("children", "op", "children commands received")
	statDelete   = Stats.Counter("delete", "op", "delete commands received")
)
