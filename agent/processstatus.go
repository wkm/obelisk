package agent

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"strings"
)

type ProcessStatus struct {
	sigar.ProcArgs
	sigar.ProcExe
	sigar.ProcMem
	sigar.ProcState
}

func errorRollup(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// get the status of all running processes
func CurrentProcessStatus() ([]ProcessStatus, error) {
	statMeasurements.Incr()

	pids := sigar.ProcList{}
	err := pids.Get()
	if err != nil {
		return nil, err
	}

	status := make([]ProcessStatus, len(pids.List))

	args := sigar.ProcArgs{}
	exe := sigar.ProcExe{}
	mem := sigar.ProcMem{}
	state := sigar.ProcState{}
	for i, pid := range pids.List {
		err := errorRollup(args.Get(pid), exe.Get(pid), mem.Get(pid), state.Get(pid))
		if err != nil {
			continue
		}
		status[i] = ProcessStatus{args, exe, mem, state}
	}

	return status, err
}

func (st *ProcessStatus) String() string {
	return fmt.Sprintf("%5d %50s  mem: %5d", st.Ppid, strings.Join(st.ProcArgs.List, " "), st.ProcMem.Size)
}
