package agent

import (
	"circuit/use/circuit"
	"fmt"
	"obelisk/lib/rlog"
	"obelisk/server/util"
)

var log = rlog.LogConfig.Logger("obelisk-agent")

const ServiceName = "obelisk-agent"

type start struct{}

func init() {
	println("registering...")
	circuit.RegisterFunc(start{})
	println("registered")
}

func (start) Start() {
	_, err := util.DiscoverObeliskServer()
	if err != nil {
		log.Printf("Error discovering obelisk server: %s", err)
		return
	}
	periodic()
}

func periodic() {
	// FIXME should report to server stats
	fmt.Printf("obelisk/agent ping")
}

// spawn an agent
func Spawn() (circuit.Addr, error) {
	_, addr, err := circuit.Spawn(
		"localhost",
		[]string{"/obelisk-agent"},
		start{},
	)

	log.Printf("spawned at %v %v", addr, err)
	return addr, err
}
