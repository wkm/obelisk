package main

import (
	"circuit/sys/acid"
	"circuit/use/circuit"
	"net/http"
	// "obelisk/lib/rinst"
	// rinstService "obelisk/lib/rinst/service"
	rlogService "obelisk/lib/rlog/service"
	"strings"
	"time"
)

type WorkerInfo struct {
	Path string
	ID   string

	Alive bool
	Error string

	// metrics components
	AcidStats     *acid.Stat
	PauseNsString string

	// stacktrace components
	RuntimeProfile string

	// profiling components
	CPUProfile string

	// metrics component
	HostInfo *HostInfo

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

	if root == "" {
		children, err := ChildrenTags("worker")
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		renderTemplate(req, rw, "/allworkers.html", children)
		return
	}

	query := req.URL.RawQuery

	log.Printf("querying for worker address")
	retrn := xServer.Call("GetWorkerAddr", root)
	log.Printf("received %#v", retrn)
	var addr circuit.Addr
	if retrn[0] != nil {
		addr = retrn[0].(circuit.Addr)
	}
	if retrn[1] != nil {
		err = retrn[1].(error)
	}

	if err != nil {
		respondError(rw, err.Error())
		return
	}

	if addr == nil {
		respondError(rw, "could not derive address")
		return
	}

	// get worker information
	workerInfo := new(WorkerInfo)
	xAcid, err := circuit.TryDial(addr, "acid")
	if err != nil {
		workerInfo.Alive = false
		workerInfo.Error = err.Error()
	} else {
		workerInfo.Alive = true
	}

	workerInfo.Path = root
	workerInfo.ID = addr.WorkerID().String()

	switch query {
	case "instrumentation":
		workers := []string{root}
		groups, err := workerResponseData(root)
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		workerInfo.HostInfo = &HostInfo{"", groups, workers}
		renderTemplate(req, rw, "/worker/worker_instrumentation.html", workerInfo)

	case "metrics":
		if !workerInfo.Alive {
			respondError(rw, "no statistics on unalive worker")
			return
		}

		workerInfo.AcidStats = xAcid.Call("Stat")[0].(*acid.Stat)
		workerInfo.PauseNsString = commaSeparated(workerInfo.AcidStats.PauseNs[:])
		renderTemplate(req, rw, "/worker/worker_stats.html", workerInfo)

	case "stacktrace":
		if !workerInfo.Alive {
			respondError(rw, "no statistics on unalive worker")
			return
		}
		retrn := xAcid.Call("RuntimeProfile", "goroutine", 1)
		workerInfo.RuntimeProfile = string(retrn[0].([]byte))
		renderTemplate(req, rw, "/worker/worker_stacktrace.html", workerInfo)

	case "profiling":
		if !workerInfo.Alive {
			respondError(rw, "no statistics on unalive worker")
			return
		}
		retrn := xAcid.Call("CPUProfile", 10*time.Second)
		workerInfo.CPUProfile = string(retrn[0].([]byte))
		renderTemplate(req, rw, "/worker/worker_profiling.html", workerInfo)

	case "logging":
		xRlog, err := circuit.TryDial(addr, rlogService.ServiceName)
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		if xRlog == nil {
			respondError(rw, "No Log service responded")
			return
		}

		retrn := xRlog.Call("FlushLog")
		workerInfo.Log = string(retrn[0].([]byte))
		renderTemplate(req, rw, "/worker/worker_logging.html", workerInfo)

	default:
		renderTemplate(req, rw, "/worker.html", workerInfo)
	}
}
