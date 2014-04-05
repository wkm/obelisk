package main

import (
	_ "circuit/load/cmd"
	"circuit/use/circuit"
	"log"
	"obelisk/agent"
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
