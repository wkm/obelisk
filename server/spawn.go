package server

import (
	"circuit/use/circuit"
	"log"
	"obelisk/lib/storetime"
)

const ServiceName = "obelisk-server"

type start struct{}

func init() {
	circuit.RegisterFunc(start{})
	circuit.RegisterValue(&ServerApp{})
	circuit.RegisterValue(&storetime.Point{})
}

func (start) Start() circuit.XPerm {
	server := new(ServerApp)
	server.Main()

	circuit.Listen(ServiceName, server)
	circuit.Daemonize(func() { <-(chan bool)(nil) })
	return circuit.PermRef(server)
}

// spawn a server on localhost
func Spawn() (circuit.Addr, error) {
	_, addr, err := circuit.Spawn(
		"localhost",
		[]string{"/obelisk-server"},
		start{},
	)

	log.Printf("spawned at %v %v", addr, err)

	return addr, err
}
