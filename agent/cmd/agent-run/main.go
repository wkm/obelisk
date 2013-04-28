package main

import (
	"fmt"
	"log"
	"obelisk/agent"
)

// run the agent on localhost
func main() {
	sys, err := agent.CurrentSystemStatus()
	if err != nil {
		log.Fatal(err.Error())
	}

	pids, err := agent.CurrentProcessStatus()
	if err != nil {
		log.Fatal(err.Error())
	}

	// system information
	fmt.Printf("%s", sys)
	println("\n")

	// process information
	for _, pid := range pids {
		if pid.Ppid != 0 {
			println(pid.String())
		}
	}

}
