package main

import (
	"net/http"
)

type WorkerInfo struct {
}

func workerHandler(rw http.ResponseWriter, req *http.Request) {
	var worker WorkerInfo
	renderTemplate(rw, "/worker.html", &worker)
}
