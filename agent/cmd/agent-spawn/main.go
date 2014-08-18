package main

import (
	"log"

	"github.com/wkm/obelisk/agent"
)

// spawn an obelisk agent worker on the local machine
func main() {
	log.Printf("Starting obelisk-agent")
	addr, err := agent.Spawn()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("obelisk-agent started %s %s", addr.String(), err)
}
