package service

import (
	"circuit/use/circuit"
	"log"
	"obelisk/lib/rinst"
)

const ServiceName = "remote-instrumentation"

func init() {
	circuit.RegisterValue(&rinst.Collection{})
}

// expose the instrumentation collection
func Expose(receiver *rinst.Collection) {
	log.Printf("exposing rinst service as %s", ServiceName)
	circuit.Listen(ServiceName, receiver)
}
