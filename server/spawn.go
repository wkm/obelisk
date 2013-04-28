package server

import (
	"circuit/use/circuit"
)

const ServiceName = "obelisk-server"

type start struct{}

func init() {
	circuit.RegisterFunc(start{})
	circuit.RegisterValue(&ServerApp{})
}

func (start) Start() circuit.XPerm {
	var server ServerApp
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

	return addr, err
}
