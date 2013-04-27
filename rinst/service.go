package rinst

import (
	"circuit/use/circuit"
)

const ServiceName = "remote-instrumentation"

func init() {
	circuit.RegisterValue(&Collection{})
}
