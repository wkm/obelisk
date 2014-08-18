package main

import (
	"github.com/wkm/obelisk/server"
	"log"
)

func main() {
	log.Printf("Starting obelisk-server")
	addr, err := server.Spawn()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("github.com/wkm/obelisk-server spawned %s", addr.String())
}
