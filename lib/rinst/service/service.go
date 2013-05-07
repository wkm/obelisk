package service

import (
	"circuit/use/circuit"
	"log"
	"obelisk/lib/rinst"
)

const (
	ServiceName = "remote-instrumentation"
	FlushSize   = 500
)

type Service struct {
	coll *rinst.Collection
}

func init() {
	circuit.RegisterValue(&rinst.Collection{})
	circuit.RegisterValue(&Service{})
}

// expose the instrumentation collection
func Expose(receiver *rinst.Collection) {
	log.Printf("exposing rinst service as %s", ServiceName)
	circuit.Listen(ServiceName, &Service{receiver})
}

func (s *Service) Schema() []rinst.Schema {
	return rinst.FlushSchema(s.coll, FlushSize)
}

func (s *Service) Measure() []rinst.Measurement {
	return rinst.FlushMeasurements(s.coll, FlushSize)
}
