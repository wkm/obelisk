package main

import (
	_ "circuit/load/cmd"
	"circuit/use/circuit"
	"log"
	"obelisk/server"
)

func main() {
	log.Printf("Starting obelisk-server")
	addr, err := server.Spawn()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("obelisk-server spawned %s", addr.String())
	circuit.Hang()
}
