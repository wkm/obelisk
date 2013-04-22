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
	Path          string
	ID            string
	AcidStats     *acid.Stat
	PauseNsString string
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
	workerInfo.AcidStats = x.Call("Stat")[0].(*acid.Stat)

	workerInfo.PauseNsString = commaSeparated(workerInfo.AcidStats.PauseNs[:])

	switch query {
	case "metrics":
		renderTemplate(rw, "/worker/worker_stats.html", workerInfo)

	default:
		renderTemplate(rw, "/worker.html", workerInfo)
	}
}
