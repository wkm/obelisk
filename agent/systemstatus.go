package agent

import (
	"github.com/cloudfoundry/gosigar"
)

type SystemStatus struct {
	sigar.Uptime
	sigar.LoadAverage
	sigar.Mem
}

// get the current system status
func CurrentSystemStatus() (SystemStatus, error) {
	statMeasurements.Incr()

	uptime := sigar.Uptime{}
	avg := sigar.LoadAverage{}
	mem := sigar.Mem{}

	var status SystemStatus
	err := errorRollup(uptime.Get(), avg.Get(), mem.Get())
	if err != nil {
		return status, err
	}

	status.Uptime = uptime
	status.LoadAverage = avg
	status.Mem = mem

	return status, nil
}
