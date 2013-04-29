package rconfig

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statSet    = Stats.Counter("set", "op", "set commands received")
	statGet    = Stats.Counter("get", "op", "get commands received")
	statGetall = Stats.Counter("getall", "op", "getall commands received")
)
