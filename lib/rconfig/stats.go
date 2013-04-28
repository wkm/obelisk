package rconfig

import (
	"obelisk/lib/rinst"
)

var Stats = rinst.NewCollection()

var (
	statSet    = Stats.Counter("set")
	statGet    = Stats.Counter("get")
	statGetall = Stats.Counter("getall")
)
