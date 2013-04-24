package main

import (
	"fmt"
	"log"
	agent "obelisk/agent/lib"
)

// run the agent once on the localhost
func main() {
	sys, err := agent.CurrentSystemStatus()
	if err != nil {
		log.Fatal(err.Error())
	}

	pids, err := agent.CurrentProcessStatus()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%s", sys)
	println("\n")
	for _, pid := range pids {
		if pid.Ppid != 0 {
			println(pid.String())
		}
	}
}
