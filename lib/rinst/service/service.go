package service

import (
	"circuit/use/circuit"
	"obelisk/lib/rinst"
)

const ServiceName = "remote-instrumentation"

func init() {
	circuit.RegisterValue(&rinst.Collection{})
}

// expose the instrumentation collection
func Expose(receiver rinst.Collection) {
	circuit.Listen(ServiceName, receiver)
}
