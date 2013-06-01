package agent

import (
	"circuit/use/circuit"
	"obelisk/lib/rconfig"
	// "obelisk/lib/rlog"
	"time"
)

const ServiceName = "obelisk-worker"

type WorkerApp struct{}
type WorkerInterface struct{}

func (WorkerApp) Main() {
	// circuit.Listen(rlog.ServiceName, rlog.Log)
	circuit.Listen(rconfig.ServiceName, rconfig.Config)
	circuit.RunInBack(func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				Periodic()
			}
		}
	})
}

func init() {
	circuit.RegisterFunc(WorkerApp{})
}
