package main

import (
	"net/http"
	// "text/template"
	anchorfs "circuit/use/anchorfs"
	"circuit/use/circuit"
	"log"
)

type zkResponse struct {
	Node  string
	Rev   int64
	Dirs  []string
	Files map[circuit.WorkerID]anchorfs.File
}

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path[3:]
	log.Printf("zk handler: %s", root)

	dir, err := anchorfs.OpenDir(root)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	dirs, err := dir.Dirs()
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	rev, files, err := dir.Files()
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	zk := zkResponse{root, rev, dirs, files}
	log.Printf("ZK: %s", zk)
	renderTemplate(rw, "/zk.html", &zk)
}
