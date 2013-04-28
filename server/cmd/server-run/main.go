package main

import (
	"log"
	"obelisk/lib/rinst"
	"obelisk/server"
	"os"
	"time"
)

func main() {
	log.Printf("Starting obelisk-server")

	var s server.ServerApp
	s.Main()

	log.Printf("Started")

	for {
		time.Sleep(5 * time.Minute)
		println("\n----")
		rinst.TextReport(os.Stdout, server.Stats)
	}

	<-(make(chan byte))
}