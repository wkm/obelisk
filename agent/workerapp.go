package agent

import (
	"circuit/use/circuit"
	"obelisk/rconfig"
	"obelisk/rlog"
	"time"
)

const ServiceName = "obelisk-worker"

type WorkerApp struct{}
type WorkerInterface struct{}

func (WorkerApp) Main() {
	circuit.Listen(rlog.ServiceName, rlog.Log)
	circuit.Listen(rconfig.ServiceName, rconfig.Config)
	circuit.Daemonize(func() {
		// this is a passive worker
		// <-(chan struct{})(nil) XXXX doesn't work
		for {
			time.Sleep(5 * time.Second)
			rlog.Log.Printf("it is now %s\n", time.Now().Format(time.Kitchen))
		}
	})
}

func init() {
	circuit.RegisterFunc(WorkerApp{})
}

func (WorkerInterface) CurrentSystemStatus() (SystemStatus, error) {
	return CurrentSystemStatus()
}

func (WorkerInterface) CurrentProcessStatus() ([]ProcessStatus, error) {
	return CurrentProcessStatus()
}
