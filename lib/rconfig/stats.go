package rconfig

import (
	"obelisk/lib/rinst"
)

var Stats = make(rinst.Collection)

var (
	statSet    = Stats.Counter("set")
	statGet    = Stats.Counter("get")
	statGetall = Stats.Counter("getall")
)
