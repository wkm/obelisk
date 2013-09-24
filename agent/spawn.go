package agent

import (
	"circuit/use/circuit"
	"obelisk/lib/rinst"
	"obelisk/lib/rlog"
	"obelisk/server/util"
	"time"
)

var log = rlog.LogConfig.Logger("obelisk-agent")

const ServiceName = "obelisk-agent"

type start struct{}

func init() {
	circuit.RegisterFunc(start{})
}

func (start) Start() {
	xServer, err := util.DiscoverObeliskServer()
	if err != nil {
		log.Printf("Error discovering obelisk server: %s", err)
		return
	}

	log.Printf("registering agent")
	xServer.Call("RegisterWorker", circuit.WorkerAddr())

	log.Printf("flushing agent schema")
	schema := rinst.FlushSchema(&StatsGauge, 10) // FIXME magic number
	xServer.Call("DeclareSchemaBuffered", circuit.WorkerAddr().WorkerID().String(), schema)

	log.Printf("starting periodic")
	circuit.RunInBack(func() { periodic(xServer) })
}

func periodic(x circuit.X) {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker
		measurements := rinst.FlushMeasurements(&StatsGauge, 10) // FIXME magic number
		x.Call("ReceiveStatsBuffered", circuit.WorkerAddr().WorkerID().String(), measurements)
	}
}

// spawn an agent
func Spawn() (circuit.Addr, error) {
	_, addr, err := circuit.Spawn(
		"localhost",
		[]string{"/obelisk-agent"},
		start{},
	)

	log.Printf("spawned at %v [err=%v]", addr, err)
	return addr, err
}
