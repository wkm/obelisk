package agent

import (
	// "time"
	// "github.com/wkm/obelisk/lib/rinst"
	"github.com/wkm/obelisk/lib/rlog"
)

var log = rlog.LogConfig.Logger("obelisk-agent")

const ServiceName = "obelisk-agent"

type start struct{}

func (start) Start() {
	// xServer, err := util.DiscoverObeliskServer()
	// if err != nil {
	// 	log.Printf("Error discovering obelisk server: %s", err)
	// 	return
	// }

	// log.Printf("registering agent")
	// xServer.Call("RegisterWorker", circuit.WorkerAddr())

	// log.Printf("flushing agent schema")
	// schema := rinst.FlushSchema(&StatsGauge, 10) // FIXME magic number
	// xServer.Call("DeclareSchemaBuffered", circuit.WorkerAddr().WorkerID().String(), schema)

	// log.Printf("starting periodic")
	// circuit.RunInBack(func() { periodic(xServer) })
}

func periodic() {
	// Flush metrics periodically
	// ticker := time.Tick(1 * time.Second)
	// for {
	// 	<-ticker
	// 	measurements := rinst.FlushMeasurements(&StatsGauge, 10) // FIXME magic number
	// 	x.Call("ReceiveStatsBuffered", circuit.WorkerAddr().WorkerID().String(), measurements)
	// }
}
