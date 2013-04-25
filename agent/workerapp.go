package agent

import (
	"circuit/use/circuit"
)

const ServiceName = "obelisk-worker"

type WorkerApp struct{}
type WorkerInterface struct{}

func (WorkerApp) Main() WorkerInterface {
	return WorkerInterface{}
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
