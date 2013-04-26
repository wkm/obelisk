package rlog

import (
	"circuit/use/circuit"
)

const ServiceName = "remote-log"

// a global object containing the log
var Log = new(Log)

func init() {
	circuit.RegisterValue(&Log{})
}
