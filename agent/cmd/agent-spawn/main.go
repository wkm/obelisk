package main

import (
	_ "circuit/kit/debug/ctrlc"
	_ "circuit/kit/debug/kill"
	_ "circuit/load"
	"circuit/use/circuit"
	"log"
	_ "net/http/pprof"
	"obelisk/agent"
)

func init() {
	println("hi")
}

// spawn an obelisk agent worker on the local machine
func main() {
	log.Printf("Starting obelisk-agent")

	_, addr, err := circuit.Spawn(
		"localhost",
		[]string{"/obelisk-agent"},
		agent.WorkerApp{},
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("obelisk-agent started %s", addr.String())
}
