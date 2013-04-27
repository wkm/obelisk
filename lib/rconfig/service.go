package rconfig

import (
	"circuit/use/circuit"
)

const ServiceName = "remote-config"

var Config = new(RConfig)

func init() {
	circuit.RegisterValue(&RConfig{})
}
