package agent

import (
	"circuit/use/circuit"
	"time"
)

const ServiceName = "obelisk-worker"

type WorkerApp struct{}
type WorkerInterface struct{}

func (WorkerApp) Main() {
	circuit.Daemonize(func() {
		time.Sleep(30 * time.Minute)
		// this is a passive worker
		<-(chan struct{})(nil)
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
