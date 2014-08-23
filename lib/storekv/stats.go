package storekv

import (
	"github.com/wkm/obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statGet = Stats.Counter("get", "op", "get commands received")
	statSet = Stats.Counter("set", "op", "set commands received")
)
