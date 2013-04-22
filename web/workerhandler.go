package main

import (
	// "circuit/sys/transport"
	"circuit/sys/acid"
	"circuit/use/circuit"
	"log"
	"net/http"
	"strings"
)

type WorkerInfo struct {
	Path string
	ID   string

	// metrics components
	AcidStats     *acid.Stat
	PauseNsString string

	// stacktrace components
	RuntimeProfile string
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
	x, err := circuit.TryDial(afile.Addr, "acid")
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	workerInfo := new(WorkerInfo)
	log.Printf("cross: %s", x)
	workerInfo.Path = root
	workerInfo.ID = afile.Addr.WorkerID().String()

	switch query {
	case "metrics":
		workerInfo.AcidStats = x.Call("Stat")[0].(*acid.Stat)
		workerInfo.PauseNsString = commaSeparated(workerInfo.AcidStats.PauseNs[:])
		renderTemplate(rw, "/worker/worker_stats.html", workerInfo)

	case "stacktrace":
		retrn := x.Call("RuntimeProfile", "goroutine", 1)
		workerInfo.RuntimeProfile = string(retrn[0].([]byte))
		renderTemplate(rw, "/worker/worker_stacktrace.html", workerInfo)

	case "profiling":
		renderTemplate(rw, "/worker/worker_profiling.html", workerInfo)

	case "logging":
		renderTemplate(rw, "/worker/worker_logging.html", workerInfo)

	default:
		renderTemplate(rw, "/worker.html", workerInfo)
	}
}
