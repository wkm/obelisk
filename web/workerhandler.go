package main

import (
	"circuit/sys/acid"
	"circuit/use/circuit"
	"net/http"
	"obelisk/rlog"
	"strings"
	"time"
)

type WorkerInfo struct {
	Path string
	ID   string

	// metrics components
	AcidStats     *acid.Stat
	PauseNsString string

	// stacktrace components
	RuntimeProfile string

	// profiling components
	CPUProfile string

	// logging bits
	Log string
}

type WorkerAddr string

func (w WorkerAddr) String() string {
	return string(w)
}

func (w WorkerAddr) Host() string {
	return strings.Split(w.String(), "@")[1]
}

func (w WorkerAddr) WorkerID() circuit.WorkerID {
	id, _ := circuit.ParseWorkerID(strings.Split(w.String(), "@")[0])
	return id
}

func workerHandler(rw http.ResponseWriter, req *http.Request) {
	// clean up URL
	root := req.URL.Path
	root = strings.TrimPrefix(root, "/worker/")
	root = strings.TrimSuffix(root, "/")

	query := req.URL.RawQuery

	data, _, err := zk.Get("/" + root)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	afile, err := getAnchorFile(data)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	// get worker information
	xAcid, err := circuit.TryDial(afile.Addr, "acid")
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	workerInfo := new(WorkerInfo)
	workerInfo.Path = root
	workerInfo.ID = afile.Addr.WorkerID().String()

	switch query {
	case "metrics":
		workerInfo.AcidStats = xAcid.Call("Stat")[0].(*acid.Stat)
		workerInfo.PauseNsString = commaSeparated(workerInfo.AcidStats.PauseNs[:])
		renderTemplate(req, rw, "/worker/worker_stats.html", workerInfo)

	case "stacktrace":
		retrn := xAcid.Call("RuntimeProfile", "goroutine", 1)
		workerInfo.RuntimeProfile = string(retrn[0].([]byte))
		renderTemplate(req, rw, "/worker/worker_stacktrace.html", workerInfo)

	case "profiling":
		retrn := xAcid.Call("CPUProfile", 10*time.Second)
		workerInfo.CPUProfile = string(retrn[0].([]byte))
		renderTemplate(req, rw, "/worker/worker_profiling.html", workerInfo)

	case "logging":
		xRlog, err := circuit.TryDial(afile.Addr, rlog.ServiceName)
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		if xRlog == nil {
			respondError(rw, "No RLog service responded")
			return
		}

		retrn := xRlog.Call("FlushLog")
		workerInfo.Log = string(retrn[0].([]byte))
		renderTemplate(req, rw, "/worker/worker_logging.html", workerInfo)

	default:
		renderTemplate(req, rw, "/worker.html", workerInfo)
	}
}
