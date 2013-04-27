package main

import (
	_ "circuit/load"
	"circuit/use/circuit"
	"log"
)

func main() {
	log.Printf("Starting obelisk-server")
	_, addr, err := circuit.Spawn(
		"localhost",
		[]string{"/obelisk-server"},
		server.WorkerApp{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("obelisk-server started %s", addr.String())
}
