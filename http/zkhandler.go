package main

import (
	"net/http"
	// "text/template"
	anchorfs "circuit/use/anchorfs"
	"circuit/use/circuit"
	"encoding/json"
	"log"
)

type zkResponse struct {
	Node    string
	Rev     int64
	Dirs    []string
	Files   []string
	workers map[circuit.WorkerID]anchorfs.File
}

func writeJson(rw http.ResponseWriter, obj interface{}) {
	encoder := json.NewEncoder(rw)
	err := encoder.Encode(&obj)
	if err != nil {
		log.Printf("json error: %s", err)
	}
}

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path[3:]
	log.Printf("zk handler: %s", root)

	dir, err := anchorfs.OpenDir(root)
	if err != nil {
		error(rw, err.Error())
		return
	}

	dirs, err := dir.Dirs()
	if err != nil {
		error(rw, err.Error())
		return
	}

	rev, files, err := dir.Files()
	if err != nil {
		error(rw, err.Error())
		return
	}

	nodes := make([]string, len(files))
	i := 0
	for id := range files {
		nodes[i] = id.String()
		i++
	}

	log.Printf("rev: %d", rev)
	log.Printf("nodes: %s", nodes)
	log.Printf("dirs: %s", dirs)

	writeJson(rw, zkResponse{root, rev, dirs, nodes, files})
}
