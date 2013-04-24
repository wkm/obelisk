package main

import (
	_ "circuit/load"
	"circuit/use/circuit"
	"log"
	agent "obelisk/agent/lib"
)

// utility command which spawns an obelisk agent worker on the local machine
func main() {
	log.Printf("Starting obelisk-agent")

	ret, addr, err := circuit.Spawn("localhost", []string{"/obelisk-agent"}, agent.App{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Agent Started")
}
