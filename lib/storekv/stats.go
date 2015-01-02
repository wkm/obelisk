package storekv

import (
	"github.com/wkm/obelisk/lib/rinst"
)

var (
	// Stats contains measurements on the usage of the key value store.
	Stats   = rinst.NewCollection()
	statGet = Stats.Counter("get", "op", "get commands received")
	statSet = Stats.Counter("set", "op", "set commands received")
)
