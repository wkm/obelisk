package agent

import (
	_ "circuit/load/cmd"
	"circuit/use/circuit"
	"obelisk/lib/rlog"
	"obelisk/server/util"
)

var log = rlog.LogConfig.Logger("obelisk-agent")

const ServiceName = "obelisk-agent"

type start struct{}

func init() {
	circuit.RegisterFunc(start{})
}

func (start) Start() {
	xServer, err := util.DiscoverObeliskServer()
	periodic()
}

func periodic() {
	// fixme
}
